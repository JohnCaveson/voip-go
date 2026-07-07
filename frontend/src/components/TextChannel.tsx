import { useState, useRef, useEffect } from 'react'
import type { Channel } from '../App'

type Message = {
  id: string
  username: string
  content: string
  created_at: string
}

type TextChannelProps = {
  channel: Channel
}

function TextChannel({ channel }: TextChannelProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = () => {
    if (!input.trim()) return

    const msg: Message = {
      id: `msg-${Date.now()}`,
      username: 'You',
      content: input.trim(),
      created_at: new Date().toISOString(),
    }

    setMessages(prev => [...prev, msg])
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
      <div className="channel-header">
        <span className="channel-icon">#</span>
        {channel.name.replace('#', '')}
      </div>

      <div className="messages">
        {messages.length === 0 && (
          <div style={{ color: '#6c7086', textAlign: 'center', marginTop: 40 }}>
            Start of #{channel.name.replace('#', '')}
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
          placeholder={`Message #${channel.name.replace('#', '')}`}
          value={input}
          onChange={e => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
        />
      </div>
    </>
  )
}

export default TextChannel
