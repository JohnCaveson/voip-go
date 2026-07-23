import { useState, useRef, useEffect } from 'react'
import type { Channel } from '../App'
import { useSignaling } from '../hooks/useSignaling'

type Message = {
  id: string
  username: string
  content: string
  created_at: string
}

type TextChannelProps = {
  channel: Channel
  signalingURL: string
  username: string
}

function TextChannel({ channel, signalingURL, username }: TextChannelProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const userIdRef = useRef('')

  const { sendTextMessage, userId } = useSignaling({
    serverUrl: signalingURL,
    room: channel.id,
    username,
    onMessage: (msg) => {
      if (msg.channelId !== channel.id) return
      const isOwn = msg.senderId === userIdRef.current
      setMessages(prev => [...prev, {
        id: `msg-${Date.now()}-${Math.random()}`,
        username: isOwn ? username : msg.senderId,
        content: msg.content,
        created_at: new Date().toISOString(),
      }])
    },
  })

  useEffect(() => {
    if (userId) userIdRef.current = userId
  }, [userId])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  useEffect(() => {
    setMessages([])
  }, [channel.id])

  const handleSend = () => {
    if (!input.trim()) return

    sendTextMessage(channel.id, input.trim())

    setMessages(prev => [...prev, {
      id: `msg-${Date.now()}-${Math.random()}`,
      username,
      content: input.trim(),
      created_at: new Date().toISOString(),
    }])
    setInput('')
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <>
      <div className="room-header">
        <span className="room-icon">💬</span>
        {channel.name}
        <span className="ephemeral-note">Messages are not persisted</span>
      </div>

      <div className="messages">
        {messages.length === 0 && (
          <div style={{ color: '#8a7e74', textAlign: 'center', marginTop: 40 }}>
            Welcome to {channel.name}
          </div>
        )}
        {messages.map(msg => (
          <div key={msg.id} className="message">
            <span className="message-username">{msg.username}</span>
            <span className="message-timestamp">
              {new Date(msg.created_at).toLocaleTimeString()}
            </span>
            <div className="message-content">{msg.content}</div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <div className="message-input-area">
        <input
          className="message-input"
          placeholder="Send a message..."
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
        />
      </div>
    </>
  )
}

export default TextChannel
