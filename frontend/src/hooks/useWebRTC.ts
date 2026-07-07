import { useRef, useCallback, useState } from 'react'

type PeerConnection = {
  pc: RTCPeerConnection
  targetId: string
}

type WebRTCHooks = {
  localStream: MediaStream | null
  isMicMuted: boolean
  isScreenSharing: boolean
  startLocalAudio: () => Promise<MediaStream>
  stopLocalAudio: () => void
  toggleMute: () => void
  startScreenShare: () => Promise<void>
  stopScreenShare: () => void
  createOffer: (targetId: string) => Promise<void>
  handleOffer: (senderId: string, sdp: string) => Promise<void>
  handleAnswer: (senderId: string, sdp: string) => Promise<void>
  handleICE: (senderId: string, candidate: string) => Promise<void>
}

export function useWebRTC(): WebRTCHooks {
  const [localStream, setLocalStream] = useState<MediaStream | null>(null)
  const [isMicMuted, setIsMicMuted] = useState(false)
  const [isScreenSharing, setIsScreenSharing] = useState(false)
  const peersRef = useRef<Map<string, PeerConnection>>(new Map())
  const screenStreamRef = useRef<MediaStream | null>(null)
  const onIceCandidateRef = useRef<((targetId: string, candidate: string) => void) | null>(null)
  const onOfferRef = useRef<((targetId: string, sdp: string) => void) | null>(null)
  const onAnswerRef = useRef<((targetId: string, sdp: string) => void) | null>(null)

  const config: RTCConfiguration = {
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' },
    ],
  }

  const onIceCandidate = useCallback((cb: (targetId: string, candidate: string) => void) => {
    onIceCandidateRef.current = cb
  }, [])

  const onOffer = useCallback((cb: (targetId: string, sdp: string) => void) => {
    onOfferRef.current = cb
  }, [])

  const onAnswer = useCallback((cb: (targetId: string, sdp: string) => void) => {
    onAnswerRef.current = cb
  }, [])

  const startLocalAudio = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      setLocalStream(stream)
      return stream
    } catch (err) {
      console.error('Failed to get microphone:', err)
      throw err
    }
  }, [])

  const stopLocalAudio = useCallback(() => {
    if (localStream) {
      localStream.getTracks().forEach(t => t.stop())
      setLocalStream(null)
    }
  }, [localStream])

  const toggleMute = useCallback(() => {
    if (localStream) {
      localStream.getAudioTracks().forEach(t => {
        t.enabled = isMicMuted
      })
    }
    setIsMicMuted(prev => !prev)
  }, [localStream, isMicMuted])

  const startScreenShare = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getDisplayMedia()
      screenStreamRef.current = stream

      stream.getVideoTracks()[0].onended = () => {
        setIsScreenSharing(false)
      }

      peersRef.current.forEach(({ pc }) => {
        stream.getVideoTracks().forEach(track => {
          pc.addTrack(track, stream)
        })
      })

      setIsScreenSharing(true)
    } catch (err) {
      console.error('Screen share failed:', err)
    }
  }, [])

  const stopScreenShare = useCallback(() => {
    if (screenStreamRef.current) {
      screenStreamRef.current.getTracks().forEach(t => t.stop())
      screenStreamRef.current = null
    }
    setIsScreenSharing(false)
  }, [])

  const createPeerConnection = useCallback(async (targetId: string): Promise<RTCPeerConnection> => {
    const pc = new RTCPeerConnection(config)

    if (localStream) {
      localStream.getTracks().forEach(track => {
        pc.addTrack(track, localStream)
      })
    }

    pc.onicecandidate = (event) => {
      if (event.candidate && onIceCandidateRef.current) {
        onIceCandidateRef.current(targetId, JSON.stringify(event.candidate))
      }
    }

    pc.ontrack = (event) => {
      const audio = document.createElement('audio')
      audio.srcObject = event.streams[0]
      audio.autoplay = true
      audio.id = `remote-${targetId}`
      document.body.appendChild(audio)

      event.streams[0].onremovetrack = () => {
        audio.remove()
      }
    }

    peersRef.current.set(targetId, { pc, targetId })
    return pc
  }, [localStream, config])

  const createOffer = useCallback(async (targetId: string) => {
    const pc = await createPeerConnection(targetId)
    const offer = await pc.createOffer()
    await pc.setLocalDescription(offer)

    if (onOfferRef.current && pc.localDescription) {
      onOfferRef.current(targetId, JSON.stringify(pc.localDescription))
    }
  }, [createPeerConnection])

  const handleOffer = useCallback(async (senderId: string, sdp: string) => {
    const pc = await createPeerConnection(senderId)
    const desc = JSON.parse(sdp) as RTCSessionDescriptionInit
    await pc.setRemoteDescription(new RTCSessionDescription(desc))

    const answer = await pc.createAnswer()
    await pc.setLocalDescription(answer)

    if (onAnswerRef.current && pc.localDescription) {
      onAnswerRef.current(senderId, JSON.stringify(pc.localDescription))
    }
  }, [createPeerConnection])

  const handleAnswer = useCallback(async (senderId: string, sdp: string) => {
    const peer = peersRef.current.get(senderId)
    if (!peer) return

    const desc = JSON.parse(sdp) as RTCSessionDescriptionInit
    await peer.pc.setRemoteDescription(new RTCSessionDescription(desc))
  }, [])

  const handleICE = useCallback(async (senderId: string, candidate: string) => {
    const peer = peersRef.current.get(senderId)
    if (!peer) return

    const iceCandidate = JSON.parse(candidate) as RTCIceCandidateInit
    await peer.pc.addIceCandidate(new RTCIceCandidate(iceCandidate))
  }, [])

  return {
    localStream,
    isMicMuted,
    isScreenSharing,
    startLocalAudio,
    stopLocalAudio,
    toggleMute,
    startScreenShare,
    stopScreenShare,
    createOffer,
    handleOffer,
    handleAnswer,
    handleICE,
  }
}
