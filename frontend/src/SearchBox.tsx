import { useState, useRef, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useSearch } from './callApi'

export default function SearchBox() {
    const [query, setQuery] = useState('')
    const [isFocused, setIsFocused] = useState(false)
    const { results, loading, error } = useSearch(query)
    const containerRef = useRef<HTMLDivElement>(null)

    // 外側クリックでサジェストを閉じる
    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setIsFocused(false)
            }
        }
        document.addEventListener('mousedown', handleClickOutside)
        return () => document.removeEventListener('mousedown', handleClickOutside)
    }, [])

    const showDropdown = isFocused && query.trim().length > 0

    return (
        <div className="relative ml-auto mr-4 flex items-center" ref={containerRef}>
            <div className="relative">
                <input
                    type="text"
                    className="w-64 px-4 py-1 bg-white border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="記事を検索..."
                    value={query}
                    onChange={e => setQuery(e.target.value)}
                    onFocus={() => setIsFocused(true)}
                />
            </div>
            
            {showDropdown && (
                <div className="absolute top-full right-0 mt-2 w-96 bg-white border border-gray-200 rounded-lg shadow-xl overflow-hidden z-50">
                    {loading && (
                        <div className="p-4 text-gray-500 text-sm text-center">検索中...</div>
                    )}
                    
                    {!loading && error && (
                        <div className="p-4 text-red-500 text-sm text-center">{error}</div>
                    )}
                    
                    {!loading && !error && results.length === 0 && (
                        <div className="p-4 text-gray-500 text-sm text-center">見つかりませんでした</div>
                    )}
                    
                    {!loading && !error && results.length > 0 && (
                        <ul className="max-h-[70vh] overflow-y-auto">
                            {results.map((res, i) => (
                                <li key={`${res.urlPath}-${i}`} className="border-b last:border-0 border-gray-100">
                                    <Link 
                                        to={res.urlPath} 
                                        className="block p-3 hover:bg-blue-50 transition-colors"
                                        onClick={() => setIsFocused(false)}
                                    >
                                        <div className="font-semibold text-gray-800">{res.displayName}</div>
                                        {res.snippet && (
                                            <div className="text-xs text-gray-500 mt-1 truncate">
                                                {res.snippet}
                                            </div>
                                        )}
                                    </Link>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>
            )}
        </div>
    )
}
