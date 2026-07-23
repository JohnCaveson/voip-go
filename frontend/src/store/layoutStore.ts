import { create } from 'zustand'
import { decodeLayout, encodeLayout } from '../utils/layoutCodec'

export type Panel = {
  id: string
  type:
    | 'channel-list'
    | 'text-channel'
    | 'voice-channel'
    | 'peers'
    | 'user-info'
    | 'connection-status'
    | 'header'
  x: number
  y: number
  width: number
  height: number
  minWidth: number
  minHeight: number
  zIndex: number
  title: string
  channelId?: string
  isMaximized: boolean
  isMinimized: boolean
  isClosable: boolean
}

type SavedBounds = {
  x: number
  y: number
  width: number
  height: number
}

type LayoutStore = {
  panels: Panel[]
  snapEnabled: boolean
  gridSize: number
  nextZIndex: number
  savedBounds: Map<string, SavedBounds>

  addPanel: (type: Panel['type'], opts?: Partial<Panel>) => void
  removePanel: (id: string) => void
  updatePanel: (id: string, updates: Partial<Panel>) => void
  bringToFront: (id: string) => void
  toggleSnap: () => void
  setGridSize: (size: number) => void
  toggleMaximize: (id: string) => void
  toggleMinimize: (id: string) => void
  resetLayout: () => void
  loadFromEncoded: (encoded: string) => boolean
  getEncoded: () => string
}

function getDefaultLayout(): Panel[] {
  const w = typeof window !== 'undefined' ? window.innerWidth : 1200
  const h = typeof window !== 'undefined' ? window.innerHeight : 800

  return [
    {
      id: 'header',
      type: 'header',
      x: 0,
      y: 0,
      width: w,
      height: 48,
      minWidth: 200,
      minHeight: 48,
      zIndex: 1,
      title: 'Gather',
      isMaximized: false,
      isMinimized: false,
      isClosable: false,
    },
    {
      id: 'channel-list',
      type: 'channel-list',
      x: 0,
      y: 48,
      width: 260,
      height: h - 48,
      minWidth: 180,
      minHeight: 200,
      zIndex: 2,
      title: 'Channels',
      isMaximized: false,
      isMinimized: false,
      isClosable: false,
    },
    {
      id: 'main-channel',
      type: 'text-channel',
      x: 260,
      y: 48,
      width: w - 260,
      height: h - 48 - 32,
      minWidth: 300,
      minHeight: 200,
      zIndex: 3,
      title: 'General',
      channelId: 'default-text',
      isMaximized: false,
      isMinimized: false,
      isClosable: false,
    },
    {
      id: 'connection-status',
      type: 'connection-status',
      x: 260,
      y: h - 32,
      width: w - 260,
      height: 32,
      minWidth: 150,
      minHeight: 32,
      zIndex: 4,
      title: 'Status',
      isMaximized: false,
      isMinimized: false,
      isClosable: true,
    },
  ]
}

export const useLayoutStore = create<LayoutStore>((set, get) => ({
  panels: getDefaultLayout(),
  snapEnabled: true,
  gridSize: 20,
  nextZIndex: 10,
  savedBounds: new Map(),

  addPanel: (type, opts) => {
    const state = get()
    const z = state.nextZIndex
    const w = typeof window !== 'undefined' ? window.innerWidth : 1200
    const h = typeof window !== 'undefined' ? window.innerHeight : 800

    const id = opts?.id ?? `${type}-${++_panelCounter}-${Date.now()}`
    const newPanel: Panel = {
      id,
      type,
      x: 50 + (state.panels.length % 5) * 30,
      y: 50 + (state.panels.length % 5) * 30,
      width: type === 'connection-status' ? 400 : 500,
      height: type === 'connection-status' ? 32 : 400,
      minWidth: 150,
      minHeight: type === 'connection-status' ? 32 : 200,
      zIndex: z,
      title: opts?.title ?? type,
      channelId: opts?.channelId,
      isMaximized: false,
      isMinimized: false,
      isClosable: true,
      ...opts,
    }

    set({
      panels: [...state.panels, newPanel],
      nextZIndex: z + 1,
    })
  },

  removePanel: (id) => {
    const state = get()
    const panel = state.panels.find((p) => p.id === id)
    if (panel && !panel.isClosable) return
    set({ panels: state.panels.filter((p) => p.id !== id) })
  },

  updatePanel: (id, updates) => {
    set({
      panels: get().panels.map((p) =>
        p.id === id ? { ...p, ...updates } : p,
      ),
    })
  },

  bringToFront: (id) => {
    const state = get()
    const z = state.nextZIndex
    set({
      panels: state.panels.map((p) =>
        p.id === id ? { ...p, zIndex: z } : p,
      ),
      nextZIndex: z + 1,
    })
  },

  toggleSnap: () => set({ snapEnabled: !get().snapEnabled }),

  setGridSize: (size) => set({ gridSize: size }),

  toggleMaximize: (id) => {
    const state = get()
    const panel = state.panels.find((p) => p.id === id)
    if (!panel) return

    const w = typeof window !== 'undefined' ? window.innerWidth : 1200
    const h = typeof window !== 'undefined' ? window.innerHeight : 800

    if (panel.isMaximized) {
      const saved = state.savedBounds.get(id)
      set({
        panels: state.panels.map((p) =>
          p.id === id
            ? {
                ...p,
                isMaximized: false,
                x: saved?.x ?? 0,
                y: saved?.y ?? 0,
                width: saved?.width ?? 500,
                height: saved?.height ?? 400,
              }
            : p,
        ),
      })
    } else {
      const saved = new Map(state.savedBounds)
      saved.set(id, {
        x: panel.x,
        y: panel.y,
        width: panel.width,
        height: panel.height,
      })
      set({
        panels: state.panels.map((p) =>
          p.id === id
            ? { ...p, isMaximized: true, x: 0, y: 0, width: w, height: h }
            : p,
        ),
        savedBounds: saved,
      })
    }
  },

  toggleMinimize: (id) => {
    set({
      panels: get().panels.map((p) =>
        p.id === id ? { ...p, isMinimized: !p.isMinimized } : p,
      ),
    })
  },

  resetLayout: () => {
    set({
      panels: getDefaultLayout(),
      nextZIndex: 10,
      savedBounds: new Map(),
    })
  },

  loadFromEncoded: (encoded) => {
    const decoded = decodeLayout(encoded)
    if (!decoded) return false
    set({
      panels: decoded.panels,
      snapEnabled: decoded.snapEnabled,
      gridSize: decoded.gridSize,
      nextZIndex: decoded.nextZIndex,
      savedBounds: new Map(),
    })
    return true
  },

  getEncoded: () => {
    const { panels, snapEnabled, gridSize } = get()
    return encodeLayout(panels, snapEnabled, gridSize)
  },
}))

let _panelCounter = 0
function newPanelId(type: string): string {
  return `${type}-${++_panelCounter}-${Date.now()}`
}
