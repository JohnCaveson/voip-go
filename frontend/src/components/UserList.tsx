type User = {
  id: string
  username: string
  is_online: boolean
}

type UserListProps = {
  users: User[]
}

function UserList({ users }: UserListProps) {
  return (
    <div className="user-list">
      <div className="channel-group-header">Online</div>
      {users.map(user => (
        <div key={user.id} className="user-item">
          <div className="user-online-dot" />
          {user.username}
        </div>
      ))}
      {users.length === 0 && (
        <div style={{ padding: '4px 16px', fontSize: 13, color: '#6c7086' }}>
          No users online
        </div>
      )}
    </div>
  )
}

export default UserList
