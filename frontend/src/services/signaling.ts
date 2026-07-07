export type SignalingMessage = {
  type: string
  room?: string
  sender_id?: string
  target_id?: string
  sdp?: string
  candidate?: string
  channel_id?: string
  content?: string
  error?: string
}

export type MessageHandler = (msg: SignalingMessage) => void

export class SignalingClient {
  private ws: WebSocket | null = null
  private handlers = new Map<string, MessageHandler[]>()

  connect(url: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(url)

      this.ws.onopen = () => resolve()
      this.ws.onerror = () => reject(new Error('WebSocket connection failed'))

      this.ws.onmessage = (event) => {
        try {
          const msg: SignalingMessage = JSON.parse(event.data)
          const typeHandlers = this.handlers.get(msg.type) || []
          typeHandlers.forEach(h => h(msg))

          const allHandlers = this.handlers.get('*') || []
          allHandlers.forEach(h => h(msg))
        } catch (err) {
          console.error('Failed to parse signaling message:', err)
        }
      }

      this.ws.onclose = () => {
        const closeHandlers = this.handlers.get('close') || []
        closeHandlers.forEach(h => h({ type: 'close' }))
      }
    })
  }

  disconnect() {
    this.ws?.close()
    this.ws = null
  }

  send(msg: SignalingMessage) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg))
    }
  }

  on(type: string, handler: MessageHandler) {
    const existing = this.handlers.get(type) || []
    existing.push(handler)
    this.handlers.set(type, existing)
  }

  off(type: string, handler: MessageHandler) {
    const existing = this.handlers.get(type) || []
    this.handlers.set(type, existing.filter(h => h !== handler))
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}
