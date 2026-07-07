export type Channel = {
  id: string
  name: string
  type: 'text' | 'voice'
  is_default: boolean
}

export type Message = {
  id: string
  channel_id: string
  user_id: string
  username: string
  content: string
  created_at: string
}

export type User = {
  id: string
  username: string
  is_online: boolean
}

class StorageService {
  private channels: Channel[] = []
  private messages: Map<string, Message[]> = new Map()
  private users: User[] = []

  async getChannels(): Promise<Channel[]> {
    return [...this.channels]
  }

  async addChannel(channel: Channel): Promise<void> {
    this.channels.push(channel)
  }

  async getMessages(channelId: string): Promise<Message[]> {
    return this.messages.get(channelId) || []
  }

  async addMessage(channelId: string, message: Message): Promise<void> {
    const existing = this.messages.get(channelId) || []
    existing.push(message)
    this.messages.set(channelId, existing)
  }

  async getUsers(): Promise<User[]> {
    return [...this.users]
  }

  async addUser(user: User): Promise<void> {
    this.users.push(user)
  }

  async setUserOnline(userId: string, online: boolean): Promise<void> {
    const user = this.users.find(u => u.id === userId)
    if (user) {
      user.is_online = online
    }
  }
}

export const storageService = new StorageService()
