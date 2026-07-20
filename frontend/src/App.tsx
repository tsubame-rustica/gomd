import { useEffect, useState } from 'react'

function App() {
  // Goから受け取るメッセージを入れるハコ（初期値は Loading...）
  const [message, setMessage] = useState<string>("Loading...")

  useEffect(() => {
    // Viteプロキシ経由で Go API にアクセス
    fetch('/api/hello')
      .then(res => res.json())
      .then(data => {
        // Goから送られてきた data.message をセット
        setMessage(data.message)
      })
      .catch(err => {
        console.error("Fetch error:", err)
        setMessage("APIからの取得に失敗しました")
      })
  }, [])

  return (
    <p className="text-xl font-medium">
      {message}
    </p>
  )
}

export default App