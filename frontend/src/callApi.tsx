import { useEffect, useState } from 'react'

// バックエンドの DocumentNode に対応する型
export interface DocumentNode {
    displayName: string
    urlPath: string
    order: number
    isFile: boolean
    children: DocumentNode[]
}


// バックエンドの API ベース URL (末尾のスラッシュを除去)
export const API_BASE = (import.meta.env.VITE_API_URL ?? '').replace(/\/+$/, '')

// HTML内の画像パスを API_BASE 付きの Cloud Run エンドポイントに解決する関数
export function resolveContentHtml(html: string, mdUrlPath: string): string {
    const lastSlash = mdUrlPath.lastIndexOf('/')
    const baseDir = lastSlash > 0 ? mdUrlPath.slice(0, lastSlash) : ''

    return html.replace(/<img\s+([^>]*?)src=["']([^"']+)["']([^>]*?)>/gi, (match, before, src, after) => {
        // 絶対URL (http://, https://, //) や data URI はそのまま
        if (/^(https?:|\/\/|data:)/i.test(src)) {
            return match
        }

        let resolvedPath: string
        if (src.startsWith('/api/contents/')) {
            resolvedPath = src
        } else if (src.startsWith('/')) {
            resolvedPath = `/api/contents${src}`
        } else {
            const cleanPath = `${baseDir}/${src}`.replace(/\/\.\//g, '/')
            resolvedPath = `/api/contents${cleanPath.startsWith('/') ? '' : '/'}${cleanPath}`
        }

        const fullSrc = API_BASE ? `${API_BASE}${resolvedPath}` : resolvedPath
        return `<img ${before}src="${fullSrc}"${after}>`
    })
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
        fetch(`${API_BASE}/api/contents${urlPath}`)
            .then(res => {
                if (!res.ok) {
                    throw new Error(`HTTP ${res.status}`)
                }
                return res.json()
            })
            .then((data: { contents: string }) => {
                const resolvedHtml = resolveContentHtml(data.contents, urlPath)
                contentCache.set(urlPath, resolvedHtml)
                setContent(resolvedHtml)
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
            fetch(`${API_BASE}/api/search?q=${encodeURIComponent(query)}`)
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