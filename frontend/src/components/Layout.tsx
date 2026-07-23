import { useCallback, useRef, useEffect } from 'react'
import Panel from './Panel'
import TextChannel from './TextChannel'
import VoiceChannel from './VoiceChannel'
import ConnectionStatus from './ConnectionStatus'
import {
  ChannelListContent,
  PeersContent,
  UserInfoContent,
  AppHeaderContent,
} from './Sidebar'
import { useLayoutStore } from '../store/layoutStore'
import type { Channel, Peer, AppConfig } from '../App'

type LayoutProps = {
  channels: Channel[]
  activeChannelId: string
  onSelectChannel: (id: string) => void
  onAddChannel: () => void
  config: AppConfig | null
  peers: Peer[]
  username: string
  signalingURL: string
}

function Layout({
  channels,
  activeChannelId,
  onSelectChannel,
  onAddChannel,
  config,
  peers,
  username,
  signalingURL,
}: LayoutProps) {
  const panels = useLayoutStore((s) => s.panels)
  const saveTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // Debounced save to backend + localStorage backup
  const scheduleSave = useCallback(() => {
    if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
    saveTimerRef.current = setTimeout(async () => {
      const encoded = useLayoutStore.getState().getEncoded()
      try {
        const { SaveLayout } = await import('../../wailsjs/go/main/App')
        await SaveLayout(encoded)
      } catch {
        // Wails not available (dev mode without Wails)
      }
      localStorage.setItem('gather-layout-backup', encoded)
    }, 500)
  }, [])

  // Subscribe to layout changes and save
  useEffect(() => {
    const unsub = useLayoutStore.subscribe(() => {
      scheduleSave()
    })
    return () => {
      unsub()
      if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
    }
  }, [scheduleSave])

  const handleSelectChannel = useCallback(
    (channelId: string) => {
      onSelectChannel(channelId)
      // Update the main channel panel
      const { updatePanel } = useLayoutStore.getState()
      const channel = channels.find((c) => c.id === channelId)
      if (channel) {
        updatePanel('main-channel', {
          channelId: channel.id,
          title: channel.name,
          type: channel.type === 'voice' ? 'voice-channel' : 'text-channel',
        })
      }
    },
    [onSelectChannel, channels],
  )

  const activeChannel = channels.find((c) => c.id === activeChannelId)

  const renderPanelContent = useCallback(
    (panelType: string, panelId: string) => {
      switch (panelType) {
        case 'header':
          return <AppHeaderContent appMode={config?.AppMode ?? 'p2p'} />
        case 'channel-list':
          return (
            <ChannelListContent
              channels={channels}
              activeChannelId={activeChannelId}
              onSelectChannel={handleSelectChannel}
              onAddChannel={onAddChannel}
            />
          )
        case 'text-channel': {
          const panel = useLayoutStore.getState().panels.find((p) => p.id === panelId)
          const ch = channels.find((c) => c.id === panel?.channelId) ?? activeChannel
          if (!ch || ch.type !== 'text') return null
          return (
            <TextChannel
              channel={ch}
              signalingURL={signalingURL}
              username={config?.Username ?? 'anonymous'}
            />
          )
        }
        case 'voice-channel': {
          const panel = useLayoutStore.getState().panels.find((p) => p.id === panelId)
          const ch = channels.find((c) => c.id === panel?.channelId) ?? activeChannel
          if (!ch || ch.type !== 'voice') return null
          return (
            <VoiceChannel
              channel={ch}
              signalingURL={signalingURL}
              username={username}
            />
          )
        }
        case 'peers':
          return <PeersContent peers={peers} />
        case 'user-info':
          return <UserInfoContent username={username} />
        case 'connection-status':
          return <ConnectionStatus config={config} />
        default:
          return null
      }
    },
    [
      channels,
      activeChannelId,
      activeChannel,
      config,
      peers,
      username,
      signalingURL,
      handleSelectChannel,
      onAddChannel,
    ],
  )

  return (
    <div className="layout-container">
      {panels.map((panel) => (
        <Panel key={panel.id} panel={panel}>
          {renderPanelContent(panel.type, panel.id)}
        </Panel>
      ))}
    </div>
  )
}

export default Layout
