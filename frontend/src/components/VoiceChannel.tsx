import { useState, useEffect, useRef } from 'react'
import type { Channel } from '../App'

type VoiceChannelProps = {
  channel: Channel
}

function VoiceChannel({ channel }: VoiceChannelProps) {
  const [isConnected, setIsConnected] = useState(false)
  const [isMuted, setIsMuted] = useState(false)
  const [isScreenSharing, setIsScreenSharing] = useState(false)
  const localAudioRef = useRef<HTMLAudioElement>(null)
  const streamRef = useRef<MediaStream | null>(null)

  const handleJoin = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      streamRef.current = stream
      if (localAudioRef.current) {
        localAudioRef.current.srcObject = stream
      }
      setIsConnected(true)
    } catch (err) {
      console.error('Failed to get mic access:', err)
    }
  }

  const handleLeave = () => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach(t => t.stop())
      streamRef.current = null
    }
    setIsConnected(false)
    setIsMuted(false)
  }

  const toggleMute = () => {
    if (streamRef.current) {
      streamRef.current.getAudioTracks().forEach(t => {
        t.enabled = isMuted
      })
    }
    setIsMuted(!isMuted)
  }

  const handleScreenShare = async () => {
    if (isScreenSharing) {
      setIsScreenSharing(false)
      return
    }

    try {
      await navigator.mediaDevices.getDisplayMedia()
      setIsScreenSharing(true)
    } catch (err) {
      console.error('Screen share failed:', err)
    }
  }

  useEffect(() => {
    return () => {
      if (streamRef.current) {
        streamRef.current.getTracks().forEach(t => t.stop())
      }
    }
  }, [])

  return (
    <>
      <div className="channel-header">
        <span className="channel-icon">🔊</span>
        {channel.name.replace('🔊 ', '')}
      </div>

      <div className="voice-controls">
        {!isConnected ? (
          <button className="voice-btn join" onClick={handleJoin}>
            Join Voice Channel
          </button>
        ) : (
          <>
            <div className="voice-users">
              <div className="voice-user">
                <div className={`speaking-indicator ${!isMuted ? 'active' : ''}`} />
                <span>You</span>
                {isMuted && <span style={{ color: '#f38ba8', fontSize: 12 }}>(muted)</span>}
              </div>
            </div>

            <div style={{ display: 'flex', gap: 8 }}>
              <button className="voice-btn leave" onClick={handleLeave}>
                Leave
              </button>
              <button className="voice-btn join" onClick={toggleMute}>
                {isMuted ? 'Unmute' : 'Mute'}
              </button>
            </div>

            <div className="screen-share-area">
              <button
                className={`screen-share-btn ${isScreenSharing ? 'stop' : 'start'}`}
                onClick={handleScreenShare}
              >
                {isScreenSharing ? 'Stop Sharing' : 'Share Screen'}
              </button>
            </div>
          </>
        )}
      </div>

      <audio ref={localAudioRef} autoPlay muted />
    </>
  )
}

export default VoiceChannel
