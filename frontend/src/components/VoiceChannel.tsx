import { useState, useEffect, useRef, useMemo } from 'react'
import type { Channel } from '../App'
import { useSignaling } from '../hooks/useSignaling'
import { useWebRTC } from '../hooks/useWebRTC'
import { getRandomJoinPhrase } from '../phrases'

type VoiceChannelProps = {
  channel: Channel
  signalingURL: string
  username: string
}

function VoiceChannel({ channel, signalingURL, username }: VoiceChannelProps) {
  const [isConnected, setIsConnected] = useState(false)
  const [remotePeers, setRemotePeers] = useState<string[]>([])
  const joinPhrase = useMemo(() => getRandomJoinPhrase(), [channel.id])

  const {
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
  } = useWebRTC()

  const signaling = useSignaling({
    serverUrl: signalingURL,
    room: channel.id,
    username,
    onPeerJoined: (peerId) => {
      setRemotePeers(prev => [...prev, peerId])
      createOffer(peerId)
    },
    onPeerLeft: (peerId) => {
      setRemotePeers(prev => prev.filter(id => id !== peerId))
    },
    onOffer: (senderId, sdp) => handleOffer(senderId, sdp),
    onAnswer: (senderId, sdp) => handleAnswer(senderId, sdp),
    onICE: (senderId, candidate) => handleICE(senderId, candidate),
  })

  const handleJoin = async () => {
    try {
      await startLocalAudio()
      setIsConnected(true)
    } catch (err) {
      console.error('Failed to get mic access:', err)
    }
  }

  const handleLeave = () => {
    stopLocalAudio()
    if (isScreenSharing) stopScreenShare()
    setIsConnected(false)
    setRemotePeers([])
  }

  useEffect(() => {
    return () => {
      stopLocalAudio()
      if (isScreenSharing) stopScreenShare()
    }
  }, [])

  return (
    <>
      <div className="room-header">
        <span className="room-icon">🎧</span>
        {channel.name}
      </div>

      <div className="voice-controls">
        {!isConnected ? (
          <button className="voice-btn join" onClick={handleJoin}>
            {joinPhrase}
          </button>
        ) : (
          <>
            <div className="voice-users">
              <div className="voice-user">
                <div className={`speaking-indicator ${!isMicMuted ? 'active' : ''}`} />
                <span>You</span>
                {isMicMuted && <span style={{ color: '#c75c3a', fontSize: 12 }}>(muted)</span>}
              </div>
              {remotePeers.map(peerId => {
                const displayName = peerId.includes('-') ? peerId.substring(0, peerId.lastIndexOf('-')) : peerId
                return (
                  <div key={peerId} className="voice-user">
                    <div className="speaking-indicator active" />
                    <span>{displayName}</span>
                  </div>
                )
              })}
            </div>

            <div style={{ display: 'flex', gap: 10 }}>
              <button className="voice-btn leave" onClick={handleLeave}>
                Leave
              </button>
              <button className="voice-btn join" onClick={toggleMute}>
                {isMicMuted ? 'Mic On' : 'Mic Off'}
              </button>
            </div>

            <div className="screen-share-area">
              <button
                className={`screen-share-btn ${isScreenSharing ? 'stop' : 'start'}`}
                onClick={isScreenSharing ? stopScreenShare : startScreenShare}
              >
                {isScreenSharing ? 'Stop Presenting' : 'Present'}
              </button>
            </div>
          </>
        )}
      </div>

      {localStream && (
        <audio ref={(el) => { if (el) el.srcObject = localStream }} autoPlay muted />
      )}
    </>
  )
}

export default VoiceChannel
