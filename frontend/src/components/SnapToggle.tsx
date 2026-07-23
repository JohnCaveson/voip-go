import { useLayoutStore } from '../store/layoutStore'

function SnapToggle() {
  const snapEnabled = useLayoutStore((s) => s.snapEnabled)
  const toggleSnap = useLayoutStore((s) => s.toggleSnap)

  return (
    <button
      className={`snap-toggle ${snapEnabled ? 'snap-toggle--on' : ''}`}
      onClick={toggleSnap}
      title={snapEnabled ? 'Snap to grid (ON) — click to disable' : 'Free-form (ON) — click to enable snap'}
    >
      <svg
        width="16"
        height="16"
        viewBox="0 0 16 16"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        {snapEnabled ? (
          <>
            <rect x="1" y="1" width="6" height="6" rx="1" fill="currentColor" />
            <rect x="9" y="1" width="6" height="6" rx="1" fill="currentColor" />
            <rect x="1" y="9" width="6" height="6" rx="1" fill="currentColor" />
            <rect x="9" y="9" width="6" height="6" rx="1" fill="currentColor" />
          </>
        ) : (
          <>
            <rect x="2" y="2" width="5" height="5" rx="1" fill="currentColor" opacity="0.5" />
            <rect x="8" y="5" width="6" height="6" rx="1" fill="currentColor" opacity="0.5" />
            <rect x="3" y="9" width="5" height="5" rx="1" fill="currentColor" opacity="0.5" />
          </>
        )}
      </svg>
      <span>{snapEnabled ? 'Grid' : 'Free'}</span>
    </button>
  )
}

export default SnapToggle
