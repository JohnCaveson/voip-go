type ScreenShareProps = {
  isSharing: boolean
  onStart: () => void
  onStop: () => void
}

function ScreenShare({ isSharing, onStart, onStop }: ScreenShareProps) {
  return (
    <div className="screen-share-area">
      <button
        className={`screen-share-btn ${isSharing ? 'stop' : 'start'}`}
        onClick={isSharing ? onStop : onStart}
      >
        {isSharing ? 'Stop Sharing Screen' : 'Share Screen'}
      </button>
    </div>
  )
}

export default ScreenShare
