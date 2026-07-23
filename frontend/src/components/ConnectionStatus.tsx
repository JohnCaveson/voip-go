import type { AppConfig } from '../App'

type ConnectionStatusProps = {
  config: AppConfig | null
}

function ConnectionStatus({ config }: ConnectionStatusProps) {
  if (!config) return null

  const isP2P = config.AppMode === 'p2p'

  return (
    <div className="connection-status-content">
      <div className="status-item">
        <span className="status-label">Storage:</span>
        <span className={`status-value ${isP2P ? 'status-local' : 'status-remote'}`}>
          {isP2P ? 'Local SQLite' : 'MongoDB'}
        </span>
      </div>
      {!isP2P && config.ServerAddr && (
        <div className="status-item">
          <span className="status-label">Server:</span>
          <span className="status-value status-connected">{config.ServerAddr}</span>
        </div>
      )}
      <div className="status-item">
        <span className="status-label">Network:</span>
        <span className="status-value">{config.NetworkMode.toUpperCase()}</span>
      </div>
    </div>
  )
}

export default ConnectionStatus
