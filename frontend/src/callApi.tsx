import { useEffect, useState } from 'react'

// バックエンドの DocumentNode に対応する型
export interface DocumentNode {
    displayName: string
    urlPath: string
    order: number
    isFile: boolean
    children: DocumentNode[]
}


// GET /api/contents/*path でMarkdownをHTMLに変換して取得するカスタムフック
export function useFetchContent(urlPath: string) {
    const [content, setContent] = useState<string>('')
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (!urlPath) return
        setLoading(true)
        fetch(`/api/contents${urlPath}`)
            .then(res => res.json())
            .then((data: { content?: string; contents?: string }) => {
                setContent(data.contents ?? data.content ?? '')
            })
            .catch(err => {
                console.error('Failed to fetch content:', err)
                setError('記事の取得に失敗しました')
            })
            .finally(() => setLoading(false))
    }, [urlPath])

    return { content, loading, error }
}

export interface SearchResult {
    urlPath: string
    displayName: string
    snippet: string
}

export function useSearch(query: string) {
    const [results, setResults] = useState<SearchResult[]>([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (!query.trim()) {
            setResults([])
            return
        }

        setLoading(true)
        // debounce的に少し待つのは呼び出し側で制御するか、ここでsetTimeoutを使う
        const timer = setTimeout(() => {
            fetch(`/api/search?q=${encodeURIComponent(query)}`)
                .then(res => res.json())
                .then(data => {
                    setResults(data.results || [])
                })
                .catch(err => {
                    console.error('Search failed:', err)
                    setError('検索に失敗しました')
                })
                .finally(() => setLoading(false))
        }, 300) // 300msデバウンス

        return () => clearTimeout(timer)
    }, [query])

    return { results, loading, error }
}