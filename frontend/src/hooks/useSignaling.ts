import { useEffect, useRef, useCallback } from 'react'
import { SignalingClient, type SignalingMessage } from '../services/signaling'

type TextMessage = {
  senderId: string
  channelId: string
  content: string
}

type UseSignalingOptions = {
  serverUrl: string
  room: string
  username: string
  onPeerJoined?: (userId: string) => void
  onPeerLeft?: (userId: string) => void
  onOffer?: (senderId: string, sdp: string) => void
  onAnswer?: (senderId: string, sdp: string) => void
  onICE?: (senderId: string, candidate: string) => void
  onMessage?: (msg: TextMessage) => void
}

export function useSignaling(options: UseSignalingOptions) {
  const clientRef = useRef<SignalingClient | null>(null)
  const userIdRef = useRef<string>('')
  const optionsRef = useRef(options)
  optionsRef.current = options

  useEffect(() => {
    if (!options.serverUrl) return

    const client = new SignalingClient()
    clientRef.current = client

    client.on('offer', (msg) => {
      if (msg.sdp && msg.sender_id) {
        optionsRef.current.onOffer?.(msg.sender_id, msg.sdp)
      }
    })

    client.on('answer', (msg) => {
      if (msg.sdp && msg.sender_id) {
        optionsRef.current.onAnswer?.(msg.sender_id, msg.sdp)
      }
    })

    client.on('ice_candidate', (msg) => {
      if (msg.candidate && msg.sender_id) {
        optionsRef.current.onICE?.(msg.sender_id, msg.candidate)
      }
    })

    client.on('peer_joined', (msg) => {
      if (msg.sender_id) {
        optionsRef.current.onPeerJoined?.(msg.sender_id)
      }
    })

    client.on('peer_left', (msg) => {
      if (msg.sender_id) {
        optionsRef.current.onPeerLeft?.(msg.sender_id)
      }
    })

    client.on('text_message', (msg) => {
      if (msg.sender_id && msg.content) {
        optionsRef.current.onMessage?.({
          senderId: msg.sender_id,
          channelId: msg.channel_id || '',
          content: msg.content,
        })
      }
    })

    const connect = async () => {
      try {
        await client.connect(options.serverUrl)
        userIdRef.current = `${options.username}-${Date.now()}`

        client.send({
          type: 'join',
          room: options.room,
          sender_id: userIdRef.current,
        })
      } catch (err) {
        console.error('Signaling connection failed:', err)
      }
    }

    connect()

    return () => {
      if (clientRef.current) {
        clientRef.current.send({
          type: 'leave',
          room: options.room,
          sender_id: userIdRef.current,
        })
        clientRef.current.disconnect()
      }
    }
  }, [options.serverUrl, options.room])

  const sendOffer = useCallback((targetId: string, sdp: string) => {
    clientRef.current?.send({
      type: 'offer',
      target_id: targetId,
      sender_id: userIdRef.current,
      sdp,
    })
  }, [])

  const sendAnswer = useCallback((targetId: string, sdp: string) => {
    clientRef.current?.send({
      type: 'answer',
      target_id: targetId,
      sender_id: userIdRef.current,
      sdp,
    })
  }, [])

  const sendICE = useCallback((targetId: string, candidate: string) => {
    clientRef.current?.send({
      type: 'ice_candidate',
      target_id: targetId,
      sender_id: userIdRef.current,
      candidate,
    })
  }, [])

  const sendTextMessage = useCallback((channelId: string, content: string) => {
    clientRef.current?.send({
      type: 'text_message',
      sender_id: userIdRef.current,
      channel_id: channelId,
      content,
    })
  }, [])

  return {
    userId: userIdRef.current,
    sendOffer,
    sendAnswer,
    sendICE,
    sendTextMessage,
    isConnected: () => clientRef.current?.isConnected() ?? false,
  }
}
