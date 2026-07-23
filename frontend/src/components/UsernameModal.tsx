import { useState } from 'react'

type UsernameModalProps = {
  onConfirm: (username: string) => void
}

function UsernameModal({ onConfirm }: UsernameModalProps) {
  const [name, setName] = useState('')

  const handleConfirm = () => {
    const trimmed = name.trim()
    if (trimmed.length > 0 && trimmed.length <= 20 && /^[a-zA-Z0-9 ]+$/.test(trimmed)) {
      onConfirm(trimmed)
    } else if (trimmed.length === 0) {
      onConfirm('anonymous')
    }
  }

  return (
    <div className="modal-overlay">
      <div className="modal username-modal">
        <div className="username-modal-header">
          <span style={{ fontSize: 32 }}>🔥</span>
          <h3>Welcome to Gather</h3>
          <p style={{ fontSize: 13, color: '#8a7e74', marginTop: 2 }}>Pick a name so others can find you</p>
        </div>

        <div>
          <label>Username</label>
          <input
            type="text"
            placeholder="e.g. Campfire Fan"
            value={name}
            onChange={e => setName(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleConfirm()}
            maxLength={20}
            autoFocus
          />
          <p style={{ fontSize: 11, color: '#8a7e74', marginTop: 4 }}>
            1–20 characters, letters, numbers, and spaces
          </p>
        </div>

        <div className="modal-actions">
          <button className="confirm" onClick={handleConfirm}>Join Gather</button>
        </div>
      </div>
    </div>
  )
}

export default UsernameModal
