import { useState, useEffect } from 'react'
import './App.css'
import Sidebar from './components/Sidebar'
import TextChannel from './components/TextChannel'
import VoiceChannel from './components/VoiceChannel'
import AddChannelModal from './components/AddChannelModal'
import UsernameModal from './components/UsernameModal'
import ConnectionStatus from './components/ConnectionStatus'
import { GetConfig, GetChannels, GetDiscoveredPeers, GetSignalingURL } from '../wailsjs/go/main/App'

export type Channel = {
  id: string
  name: string
  type: 'text' | 'voice'
  is_default: boolean
}

export type Peer = {
  id: string
  username: string
  addr: string
  signaling_addr: string
}

export type AppConfig = {
  AppMode: string
  NetworkMode: string
  ServerAddr: string
  MongoDBURI: string
  Username: string
}

function App() {
  const [channels, setChannels] = useState<Channel[]>([])
  const [activeChannelId, setActiveChannelId] = useState('default-text')
  const [showAddModal, setShowAddModal] = useState(false)
  const [config, setConfig] = useState<AppConfig | null>(null)
  const [peers, setPeers] = useState<Peer[]>([])
  const [signalingURL, setSignalingURL] = useState('')
  const [username, setUsername] = useState<string | null>(() => localStorage.getItem('gather-username'))
  const [showUsernameModal, setShowUsernameModal] = useState(false)

  useEffect(() => {
    loadConfig()
    loadChannels()
    loadSignalingURL()
  }, [])

  useEffect(() => {
    if (username === null) {
      setShowUsernameModal(true)
    }
  }, [username])

  useEffect(() => {
    if (config?.AppMode !== 'p2p') return
    const interval = setInterval(loadPeers, 5000)
    loadPeers()
    return () => clearInterval(interval)
  }, [config])

  const loadConfig = async () => {
    try {
      const cfg = await GetConfig() as AppConfig
      setConfig(cfg)
    } catch (err) {
      console.error('Failed to load config:', err)
      setConfig({ AppMode: 'p2p', NetworkMode: 'lan', ServerAddr: '', MongoDBURI: '', Username: 'anonymous' })
    }
  }

  const loadChannels = async () => {
    try {
      const chs = await GetChannels()
      if (chs && chs.length > 0) {
        const mapped: Channel[] = chs.map(ch => ({
          id: ch.id,
          name: ch.name,
          type: ch.type as 'text' | 'voice',
          is_default: ch.is_default,
        }))
        setChannels(mapped)
        if (!activeChannelId || !mapped.find(c => c.id === activeChannelId)) {
          setActiveChannelId(mapped[0].id)
        }
      } else {
        setChannels([
          { id: 'default-text', name: 'General', type: 'text', is_default: true },
          { id: 'default-voice', name: 'Lounge', type: 'voice', is_default: true },
        ])
      }
    } catch (err) {
      console.error('Failed to load channels:', err)
      setChannels([
        { id: 'default-text', name: 'General', type: 'text', is_default: true },
        { id: 'default-voice', name: 'Lounge', type: 'voice', is_default: true },
      ])
    }
  }

  const loadPeers = async () => {
    try {
      const p = await GetDiscoveredPeers()
      setPeers(p || [])
    } catch (err) {
      console.error('Failed to load peers:', err)
    }
  }

  const loadSignalingURL = async () => {
    try {
      const url = await GetSignalingURL()
      setSignalingURL(url)
    } catch (err) {
      console.error('Failed to get signaling URL:', err)
    }
  }

  const activeChannel = channels.find(c => c.id === activeChannelId)

  const handleSetUsername = (name: string) => {
    localStorage.setItem('gather-username', name)
    setUsername(name)
    setShowUsernameModal(false)
  }

  const handleAddChannel = async (name: string, type: 'text' | 'voice') => {
    try {
      await GetChannels()
      await loadChannels()
    } catch (err) {
      const newChannel: Channel = {
        id: `ch-${Date.now()}`,
        name: name,
        type,
        is_default: false,
      }
      setChannels(prev => [...prev, newChannel])
      setActiveChannelId(newChannel.id)
    }
  }

  return (
    <div className="app">
      <Sidebar
        channels={channels}
        activeChannelId={activeChannelId}
        onSelectChannel={setActiveChannelId}
        onAddChannel={() => setShowAddModal(true)}
        appMode={config?.AppMode ?? 'p2p'}
        peers={peers}
        username={username ?? 'anonymous'}
      />
      <div className="main-content">
        {activeChannel?.type === 'text' && (
          <TextChannel
            channel={activeChannel}
            signalingURL={signalingURL}
            username={config?.Username ?? 'anonymous'}
          />
        )}
        {activeChannel?.type === 'voice' && (
          <VoiceChannel
            channel={activeChannel}
            signalingURL={signalingURL}
            username={username ?? 'anonymous'}
          />
        )}
        <ConnectionStatus config={config} />
      </div>
      {showAddModal && (
        <AddChannelModal
          onClose={() => setShowAddModal(false)}
          onConfirm={handleAddChannel}
        />
      )}
      {showUsernameModal && (
        <UsernameModal onConfirm={handleSetUsername} />
      )}
    </div>
  )
}

export default App
