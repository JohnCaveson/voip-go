import { useState } from 'react'

type AddChannelModalProps = {
  onClose: () => void
  onConfirm: (name: string, type: 'text' | 'voice') => void
}

function AddChannelModal({ onClose, onConfirm }: AddChannelModalProps) {
  const [name, setName] = useState('')
  const [type, setType] = useState<'text' | 'voice'>('text')

  const handleConfirm = () => {
    if (!name.trim()) return
    onConfirm(name.trim(), type)
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" onClick={e => e.stopPropagation()}>
        <h3>New Room</h3>

        <div>
          <label>Room Type</label>
          <select value={type} onChange={e => setType(e.target.value as 'text' | 'voice')}>
            <option value="text">Chat</option>
            <option value="voice">Audio</option>
          </select>
        </div>

        <div>
          <label>Room Name</label>
          <input
            type="text"
            placeholder="e.g. Game Night"
            value={name}
            onChange={e => setName(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleConfirm()}
            autoFocus
          />
        </div>

        <div className="modal-actions">
          <button className="cancel" onClick={onClose}>Cancel</button>
          <button className="confirm" onClick={handleConfirm}>Create</button>
        </div>
      </div>
    </div>
  )
}

export default AddChannelModal
