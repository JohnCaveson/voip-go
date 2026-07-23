import type { Panel } from '../store/layoutStore'

type EncodedLayout = {
  v: number
  s: boolean
  g: number
  p: EncodedPanel[]
}

type EncodedPanel = {
  i: string
  t: string
  x: number
  y: number
  w: number
  h: number
  z: number
  c: boolean
  ch?: string
  ti?: string
}

function encodePanel(p: Panel): EncodedPanel {
  const ep: EncodedPanel = {
    i: p.id,
    t: p.type,
    x: p.x,
    y: p.y,
    w: p.width,
    h: p.height,
    z: p.zIndex,
    c: p.isClosable,
  }
  if (p.channelId) ep.ch = p.channelId
  if (p.title) ep.ti = p.title
  return ep
}

function decodePanelType(t: string): Panel['type'] {
  const valid: Panel['type'][] = [
    'channel-list', 'text-channel', 'voice-channel',
    'peers', 'user-info', 'connection-status', 'header',
  ]
  return valid.includes(t as Panel['type']) ? (t as Panel['type']) : 'text-channel'
}

function defaultMinWidth(type: Panel['type']): number {
  switch (type) {
    case 'channel-list': return 180
    case 'connection-status': return 150
    case 'header': return 200
    case 'user-info': return 150
    default: return 300
  }
}

function defaultMinHeight(type: Panel['type']): number {
  switch (type) {
    case 'connection-status': return 30
    case 'header': return 40
    case 'user-info': return 50
    default: return 200
  }
}

function defaultTitle(type: Panel['type']): string {
  switch (type) {
    case 'header': return 'Gather'
    case 'channel-list': return 'Channels'
    case 'connection-status': return 'Status'
    case 'user-info': return 'User'
    case 'peers': return 'Peers'
    default: return 'Channel'
  }
}

export function encodeLayout(
  panels: Panel[],
  snapEnabled: boolean,
  gridSize: number,
): string {
  const data: EncodedLayout = {
    v: 1,
    s: snapEnabled,
    g: gridSize,
    p: panels.map(encodePanel),
  }
  return btoa(JSON.stringify(data))
}

export function decodeLayout(encoded: string): {
  panels: Panel[]
  snapEnabled: boolean
  gridSize: number
  nextZIndex: number
} | null {
  try {
    const data: EncodedLayout = JSON.parse(atob(encoded))
    if (data.v !== 1) return null

    const panels: Panel[] = data.p.map((ep) => ({
      id: ep.i,
      type: decodePanelType(ep.t),
      x: ep.x,
      y: ep.y,
      width: ep.w,
      height: ep.h,
      minWidth: defaultMinWidth(decodePanelType(ep.t)),
      minHeight: defaultMinHeight(decodePanelType(ep.t)),
      zIndex: ep.z,
      title: ep.ti ?? defaultTitle(decodePanelType(ep.t)),
      channelId: ep.ch,
      isMaximized: false,
      isMinimized: false,
      isClosable: ep.c,
    }))

    return {
      panels,
      snapEnabled: data.s,
      gridSize: data.g,
      nextZIndex: Math.max(0, ...panels.map((p) => p.zIndex)) + 1,
    }
  } catch {
    return null
  }
}
