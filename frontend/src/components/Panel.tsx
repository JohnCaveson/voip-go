import { useCallback, type ReactNode } from 'react'
import { Rnd } from 'react-rnd'
import { useLayoutStore, type Panel as PanelType } from '../store/layoutStore'

type PanelProps = {
  panel: PanelType
  children: ReactNode
}

function Panel({ panel, children }: PanelProps) {
  const snapEnabled = useLayoutStore((s) => s.snapEnabled)
  const gridSize = useLayoutStore((s) => s.gridSize)
  const bringToFront = useLayoutStore((s) => s.bringToFront)
  const updatePanel = useLayoutStore((s) => s.updatePanel)
  const toggleMaximize = useLayoutStore((s) => s.toggleMaximize)
  const toggleMinimize = useLayoutStore((s) => s.toggleMinimize)
  const removePanel = useLayoutStore((s) => s.removePanel)

  const dragGrid: [number, number] = snapEnabled ? [gridSize, gridSize] : [1, 1]

  const handleMouseDown = useCallback(() => {
    bringToFront(panel.id)
  }, [bringToFront, panel.id])

  const handleDragStop = useCallback(
    (_: unknown, d: { x: number; y: number }) => {
      updatePanel(panel.id, { x: d.x, y: d.y })
    },
    [updatePanel, panel.id],
  )

  const handleResizeStop = useCallback(
    (
      _: unknown,
      __: unknown,
      ref: HTMLElement,
      delta: { width: number; height: number },
      position: { x: number; y: number },
    ) => {
      updatePanel(panel.id, {
        x: position.x,
        y: position.y,
        width: parseInt(ref.style.width, 10),
        height: parseInt(ref.style.height, 10),
      })
    },
    [updatePanel, panel.id],
  )

  if (panel.isMinimized) {
    return (
      <Rnd
        size={{ width: panel.width, height: 36 }}
        position={{ x: panel.x, y: panel.y }}
        disableDragging
        enableResizing={false}
        onMouseDown={handleMouseDown}
        style={{ zIndex: panel.zIndex }}
      >
        <div className="panel panel--minimized">
          <div className="panel-header panel-header--minimized">
            <span className="panel-title">{panel.title}</span>
            <div className="panel-controls">
              <button
                className="panel-btn panel-btn-restore"
                onClick={() => toggleMinimize(panel.id)}
                title="Restore"
              >
                □
              </button>
              {panel.isClosable && (
                <button
                  className="panel-btn panel-btn-close"
                  onClick={() => removePanel(panel.id)}
                  title="Close"
                >
                  ×
                </button>
              )}
            </div>
          </div>
        </div>
      </Rnd>
    )
  }

  return (
    <Rnd
      size={{ width: panel.width, height: panel.height }}
      position={{ x: panel.x, y: panel.y }}
      dragGrid={dragGrid}
      resizeGrid={dragGrid}
      bounds="parent"
      enableResizing={{
        top: true,
        right: true,
        bottom: true,
        left: true,
        topRight: true,
        bottomRight: true,
        bottomLeft: true,
        topLeft: true,
      }}
      onMouseDown={handleMouseDown}
      onDragStop={handleDragStop}
      onResizeStop={handleResizeStop}
      style={{ zIndex: panel.zIndex }}
      minWidth={panel.minWidth}
      minHeight={panel.minHeight}
    >
      <div
        className={`panel ${panel.isMaximized ? 'panel--maximized' : ''}`}
      >
        <div
          className="panel-header"
          onMouseDown={(e) => {
            if ((e.target as HTMLElement).closest('.panel-controls')) return
            handleMouseDown()
          }}
        >
          <span className="panel-title">{panel.title}</span>
          <div className="panel-controls">
            <button
              className="panel-btn panel-btn-minimize"
              onClick={() => toggleMinimize(panel.id)}
              title="Minimize"
            >
              −
            </button>
            <button
              className="panel-btn panel-btn-maximize"
              onClick={() => toggleMaximize(panel.id)}
              title={panel.isMaximized ? 'Restore' : 'Maximize'}
            >
              {panel.isMaximized ? '❐' : '□'}
            </button>
            {panel.isClosable && (
              <button
                className="panel-btn panel-btn-close"
                onClick={() => removePanel(panel.id)}
                title="Close"
              >
                ×
              </button>
            )}
          </div>
        </div>
        <div className="panel-content">{children}</div>
      </div>
    </Rnd>
  )
}

export default Panel
