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
        <h3>Add Channel</h3>

        <div>
          <label>Channel Type</label>
          <select value={type} onChange={e => setType(e.target.value as 'text' | 'voice')}>
            <option value="text">Text</option>
            <option value="voice">Voice</option>
          </select>
        </div>

        <div>
          <label>Channel Name</label>
          <input
            type="text"
            placeholder={type === 'text' ? 'channel-name' : 'Channel Name'}
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
