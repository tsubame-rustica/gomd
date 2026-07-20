import { useEffect, useState } from 'react'

function App() {
  // Goから受け取るメッセージを入れるハコ（初期値は Loading...）
  const [content, setContent] = useState<string>("Loading...")

  useEffect(() => {
    // Viteプロキシ経由で Go API にアクセス
    fetch('/api/memo/test/hoge')
      .then(res => res.json())
      .then(data => {
        // Goから送られてきた data.message をセット
        setContent(data.message)
      })
      .catch(err => {
        console.error("Fetch error:", err)
        setContent("APIからの取得に失敗しました")
      })
  }, [])

  return (
    <div className="prose prose-blue" dangerouslySetInnerHTML={{ __html: content }} />
  )
}

export default App