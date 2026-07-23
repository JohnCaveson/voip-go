import type { Channel, Peer } from '../App'

type ChannelListProps = {
  channels: Channel[]
  activeChannelId: string
  onSelectChannel: (id: string) => void
  onAddChannel: () => void
}

function ChannelListContent({ channels, activeChannelId, onSelectChannel, onAddChannel }: ChannelListProps) {
  const chatRooms = channels.filter(c => c.type === 'text')
  const audioRooms = channels.filter(c => c.type === 'voice')

  return (
    <div className="channel-list">
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
    </div>
  )
}

type PeersPanelProps = {
  peers: Peer[]
}

function PeersContent({ peers }: PeersPanelProps) {
  return (
    <div className="channel-list">
      <div className="room-group">
        {peers.length === 0 ? (
          <div style={{ color: '#8a7e74', padding: '4px 18px', fontSize: 13 }}>
            No peers found
          </div>
        ) : (
          peers.map(peer => (
            <div key={peer.id} className="room-item peer">
              <span className="peer-dot" />
              {peer.username}
            </div>
          ))
        )}
      </div>
    </div>
  )
}

type UserInfoProps = {
  username: string
}

function UserInfoContent({ username }: UserInfoProps) {
  return (
    <div className="user-info-panel">
      <div className="sidebar-user">
        <div className="user-online-dot" />
        <span>{username}</span>
      </div>
    </div>
  )
}

type AppHeaderProps = {
  appMode: string
}

function AppHeaderContent({ appMode }: AppHeaderProps) {
  const isP2P = appMode === 'p2p'

  return (
    <div className="app-header-content">
      <div className="sidebar-header-left">
        <span className="app-logo">🔥</span>
        <span>Gather</span>
      </div>
      <span className={`mode-badge ${isP2P ? 'mode-p2p' : 'mode-hosted'}`}>
        {isP2P ? 'P2P' : 'Hosted'}
      </span>
    </div>
  )
}

export { ChannelListContent, PeersContent, UserInfoContent, AppHeaderContent }
