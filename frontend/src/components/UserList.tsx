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
      <div className="room-group-header">Active</div>
      {users.map(user => (
        <div key={user.id} className="user-item">
          <div className="user-online-dot" />
          {user.username}
        </div>
      ))}
      {users.length === 0 && (
        <div style={{ padding: '4px 18px', fontSize: 13, color: '#8a7e74' }}>
          No one here yet
        </div>
      )}
    </div>
  )
}

export default UserList
