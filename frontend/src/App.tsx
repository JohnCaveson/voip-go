import { useState, useEffect } from 'react'
import './App.css'
import Sidebar from './components/Sidebar'
import TextChannel from './components/TextChannel'
import VoiceChannel from './components/VoiceChannel'
import AddChannelModal from './components/AddChannelModal'
import ConnectionStatus from './components/ConnectionStatus'
import { GetConfig, GetChannels } from '../wailsjs/go/main/App'

export type Channel = {
  id: string
  name: string
  type: 'text' | 'voice'
  is_default: boolean
}

export type User = {
  id: string
  username: string
  is_online: boolean
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

  useEffect(() => {
    loadConfig()
    loadChannels()
  }, [])

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

  const activeChannel = channels.find(c => c.id === activeChannelId)

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
      />
      <div className="main-content">
        {activeChannel?.type === 'text' && (
          <TextChannel channel={activeChannel} />
        )}
        {activeChannel?.type === 'voice' && (
          <VoiceChannel channel={activeChannel} />
        )}
        <ConnectionStatus config={config} />
      </div>
      {showAddModal && (
        <AddChannelModal
          onClose={() => setShowAddModal(false)}
          onConfirm={handleAddChannel}
        />
      )}
    </div>
  )
}

export default App
