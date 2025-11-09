import { useState, useEffect } from 'react'

type Role = 'EMS' | 'CITIZEN'

function App() {
  const [messages, setMessages] = useState<string[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [role, setRole] = useState<Role>('CITIZEN')
  const [apiBase, setApiBase] = useState<string>(window.location.origin)

  useEffect(() => {
    let mounted = true
    const fetchMessages = async () => {
      try {
        const response = await fetch(`${apiBase}/messages`)
        const data = await response.json()
        if (!mounted) return
        // Expecting an array of strings (backend stores human-readable lines)
        if (Array.isArray(data)) {
          setMessages(data)
        } else {
          // Fallback: try to coerce structured messages to strings
          const asStrings = (data || []).map((m: any) => {
            if (typeof m === 'string') return m
            if (m.content) return m.content
            return JSON.stringify(m)
          })
          setMessages(asStrings)
        }
      } catch (err) {
        console.error('Failed to fetch messages:', err)
      }
    }

    fetchMessages()
    const interval = setInterval(fetchMessages, 2000)
    return () => {
      mounted = false
      clearInterval(interval)
    }
  }, [apiBase])

  const sendMessage = async () => {
    if (!newMessage.trim()) return

    try {
      const body = { message: newMessage }
      const response = await fetch(`${apiBase}/send`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })

      if (!response.ok) throw new Error('Failed to send message')
      setNewMessage('')
    } catch (err) {
      console.error('Failed to send message:', err)
    }
  }

  return (
    <div className="min-h-screen bg-[#f8f9fa]">
      <div className="max-w-3xl mx-auto py-4 px-4">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-2xl font-bold">MeshComm</h1>
          <div className="flex items-center space-x-3">
            <label className="text-sm">Role:</label>
            <select
              value={role}
              onChange={(e) => setRole(e.target.value as Role)}
              className="px-3 py-1 border rounded"
            >
              <option value="CITIZEN">Citizen</option>
              <option value="EMS">Emergency Services</option>
            </select>
            <label className="text-sm">API:</label>
            <input
              className="px-2 py-1 border rounded w-48"
              value={apiBase}
              onChange={(e) => setApiBase(e.target.value)}
            />
          </div>
        </div>

        {/* Messages Display */}
        <div className="h-[60vh] overflow-y-auto px-2 mb-4">
          <div className="space-y-3">
            {messages.length === 0 ? (
              <div className="text-gray-500">No messages yet</div>
            ) : (
              messages.map((msg, idx) => (
                <div key={idx} className="bg-white rounded-2xl px-4 py-3 shadow-sm">
                  {msg}
                </div>
              ))
            )}
          </div>
        </div>

        {/* Message Input */}
        <div className="mt-4">
          <div className="flex space-x-4">
            <input
              type="text"
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder={role === 'EMS' ? 'Type an alert or update...' : 'Type your message...'}
              className="flex-1 px-4 py-3 border rounded-full focus:outline-none"
              onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
            />
            <button
              onClick={sendMessage}
              className="px-6 py-3 bg-[#90EE90] text-white rounded-full hover:bg-[#77dd77] transition-colors"
              disabled={!newMessage.trim()}
            >
              send
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App