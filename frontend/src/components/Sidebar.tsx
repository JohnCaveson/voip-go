import type { Channel, Peer } from '../App'

type SidebarProps = {
  channels: Channel[]
  activeChannelId: string
  onSelectChannel: (id: string) => void
  onAddChannel: () => void
  appMode: string
  peers: Peer[]
  username: string
}

function Sidebar({ channels, activeChannelId, onSelectChannel, onAddChannel, appMode, peers, username }: SidebarProps) {
  const chatRooms = channels.filter(c => c.type === 'text')
  const audioRooms = channels.filter(c => c.type === 'voice')
  const isP2P = appMode === 'p2p'

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <div className="sidebar-header-left">
          <span className="app-logo">🔥</span>
          <span>Gather</span>
        </div>
        <span className={`mode-badge ${isP2P ? 'mode-p2p' : 'mode-hosted'}`}>
          {isP2P ? 'P2P' : 'Hosted'}
        </span>
      </div>

      <div className="room-group">
        <div className="room-group-header">
          <span>Chat</span>
          <button className="room-group-add" onClick={onAddChannel} title="New room">+</button>
        </div>
        {chatRooms.map(ch => (
          <div
            key={ch.id}
            className={`room-item ${ch.id === activeChannelId ? 'active' : ''}`}
            onClick={() => onSelectChannel(ch.id)}
          >
            <span className="room-icon">💬</span>
            {ch.name}
          </div>
        ))}
      </div>

      <div className="room-group">
        <div className="room-group-header">
          <span>Audio</span>
        </div>
        {audioRooms.map(ch => (
          <div
            key={ch.id}
            className={`room-item ${ch.id === activeChannelId ? 'active' : ''}`}
            onClick={() => onSelectChannel(ch.id)}
          >
            <span className="room-icon">🎧</span>
            {ch.name}
          </div>
        ))}
      </div>

      {isP2P && peers.length > 0 && (
        <div className="room-group">
          <div className="room-group-header">
            <span>Peers Nearby</span>
          </div>
          {peers.map(peer => (
            <div key={peer.id} className="room-item peer">
              <span className="peer-dot" />
              {peer.username}
            </div>
          ))}
        </div>
      )}

      <div className="sidebar-footer">
        <div className="sidebar-user">
          <div className="user-online-dot" />
          <span>{username}</span>
        </div>
      </div>
    </div>
  )
}

export default Sidebar
