import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import Sidebar from '../Sidebar'

function Post() {
  const { category, title } = useParams<{ category: string; title: string }>()
  const [content, setContent] = useState<string>("Loading...")

  if (!category || !title) {
    return <p>Invalid URL parameters</p>
  }

  useEffect(() => {
    // Viteプロキシ経由で Go API にアクセス
    fetch(`/api/memo/getPost/${category}/${title}`)
      .then(res => res.json())
      .then(data => {
        // Goから送られてきた data.content をセット
        if (data.content) {
            setContent(data.content)
        } else {
            setContent("<p>コンテンツが見つかりませんでした</p>")
        }
      })
      .catch(err => {
        console.error("Fetch error:", err)
        setContent("<p>エラーが発生しました</p>")
      })
  }, [category, title])

  return (
    <>
        <Sidebar category={category} />
        <main className="flex-1 overflow-y-scroll">
            <div className="prose prose-blue ml-8 m-4" dangerouslySetInnerHTML={{ __html: content }} />
        </main>
    </>
  )
}

export default Post