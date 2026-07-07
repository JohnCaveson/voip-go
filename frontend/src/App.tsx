import { useState } from 'react'
import './App.css'
import Sidebar from './components/Sidebar'
import TextChannel from './components/TextChannel'
import VoiceChannel from './components/VoiceChannel'
import AddChannelModal from './components/AddChannelModal'

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

function App() {
  const [channels, setChannels] = useState<Channel[]>([
    { id: 'default-text', name: '#general', type: 'text', is_default: true },
    { id: 'default-voice', name: '🔊 General', type: 'voice', is_default: true },
  ])
  const [activeChannelId, setActiveChannelId] = useState('default-text')
  const [showAddModal, setShowAddModal] = useState(false)

  const activeChannel = channels.find(c => c.id === activeChannelId)

  const handleAddChannel = (name: string, type: 'text' | 'voice') => {
    const newChannel: Channel = {
      id: `ch-${Date.now()}`,
      name: type === 'text' ? `#${name}` : `🔊 ${name}`,
      type,
      is_default: false,
    }
    setChannels(prev => [...prev, newChannel])
    setActiveChannelId(newChannel.id)
  }

  return (
    <div className="app">
      <Sidebar
        channels={channels}
        activeChannelId={activeChannelId}
        onSelectChannel={setActiveChannelId}
        onAddChannel={() => setShowAddModal(true)}
      />
      <div className="main-content">
        {activeChannel?.type === 'text' && (
          <TextChannel channel={activeChannel} />
        )}
        {activeChannel?.type === 'voice' && (
          <VoiceChannel channel={activeChannel} />
        )}
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
