import type { Channel } from '../App'

type SidebarProps = {
  channels: Channel[]
  activeChannelId: string
  onSelectChannel: (id: string) => void
  onAddChannel: () => void
}

function Sidebar({ channels, activeChannelId, onSelectChannel, onAddChannel }: SidebarProps) {
  const textChannels = channels.filter(c => c.type === 'text')
  const voiceChannels = channels.filter(c => c.type === 'voice')

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <span>VoIP App</span>
      </div>

      <div className="channel-group">
        <div className="channel-group-header">Text Channels</div>
        {textChannels.map(ch => (
          <div
            key={ch.id}
            className={`channel-item ${ch.id === activeChannelId ? 'active' : ''}`}
            onClick={() => onSelectChannel(ch.id)}
          >
            <span className="channel-icon">#</span>
            {ch.name.replace('#', '')}
          </div>
        ))}
      </div>

      <div className="channel-group">
        <div className="channel-group-header">Voice Channels</div>
        {voiceChannels.map(ch => (
          <div
            key={ch.id}
            className={`channel-item ${ch.id === activeChannelId ? 'active' : ''}`}
            onClick={() => onSelectChannel(ch.id)}
          >
            <span className="channel-icon">🔊</span>
            {ch.name.replace('🔊 ', '')}
          </div>
        ))}
      </div>

      <button className="add-channel-btn" onClick={onAddChannel}>
        + Add Channel
      </button>
    </div>
  )
}

export default Sidebar
