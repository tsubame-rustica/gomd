import { useEffect, useState } from 'react'

// バックエンドの DocumentNode に対応する型
export interface DocumentNode {
    displayName: string
    urlPath: string
    order: number
    isFile: boolean
    children: DocumentNode[]
}


// 記事HTMLコンテンツのインメモリキャッシュ (urlPath -> HTML文字列)
const contentCache = new Map<string, string>()

// GET /api/contents/*path でMarkdownをHTMLに変換して取得するカスタムフック
export function useFetchContent(urlPath: string) {
    const [content, setContent] = useState<string>(() => contentCache.get(urlPath) ?? '')
    const [loading, setLoading] = useState<boolean>(() => !urlPath || !contentCache.has(urlPath))
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (!urlPath) {
            setContent('')
            setLoading(false)
            return
        }

        // キャッシュに既に存在する場合は fetch せずに即時反映
        const cached = contentCache.get(urlPath)
        if (cached !== undefined) {
            setContent(cached)
            setLoading(false)
            setError(null)
            return
        }

        setLoading(true)
        setError(null)
        fetch(`/api/contents${urlPath}`)
            .then(res => {
                if (!res.ok) {
                    throw new Error(`HTTP ${res.status}`)
                }
                return res.json()
            })
            .then((data: { contents: string }) => {
                contentCache.set(urlPath, data.contents)
                setContent(data.contents)
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