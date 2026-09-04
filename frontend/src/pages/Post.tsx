import { useLocation } from 'react-router-dom'
import Sidebar from '../Sidebar'
import { useFetchContent } from '../callApi'

function Post() {
    const location = useLocation()
    const urlPath = location.pathname  // 例: /git/default/default.md

    // .md で終わるパスのみコンテンツを取得する
    const isMd = urlPath.endsWith('.md')
    const { content, loading, error } = useFetchContent(isMd ? urlPath : '')

    return (
        <>
            <Sidebar />
            <main className="flex-1 overflow-y-scroll">
                {!isMd && (
                    <p className="m-8 text-neutral-400">← サイドバーから記事を選んでください</p>
                )}
                {isMd && loading && <p className="m-8 text-neutral-500">読み込み中...</p>}
                {isMd && error && <p className="m-8 text-red-500">{error}</p>}
                {isMd && !loading && !error && (
                    <div
                        className="prose prose-blue ml-8 m-4"
                        dangerouslySetInnerHTML={{ __html: content }}
                    />
                )}
            </main>
        </>
    )
}

export default Post