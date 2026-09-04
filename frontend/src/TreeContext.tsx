import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import type { DocumentNode } from './callApi'

interface TreeContextValue {
    tree: DocumentNode | null
    loading: boolean
    error: string | null
}

const TreeContext = createContext<TreeContextValue>({
    tree: null,
    loading: true,
    error: null,
})

// ツリーを一度だけ fetch して全コンポーネントへ提供するプロバイダー
export function TreeProvider({ children }: { children: ReactNode }) {
    const [tree, setTree] = useState<DocumentNode | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        fetch('/api/tree')
            .then(res => res.json())
            .then((data: DocumentNode) => setTree(data))
            .catch(err => {
                console.error('Failed to fetch tree:', err)
                setError('ツリーの取得に失敗しました')
            })
            .finally(() => setLoading(false))
    }, [])

    return (
        <TreeContext.Provider value={{ tree, loading, error }}>
            {children}
        </TreeContext.Provider>
    )
}

// Header / Sidebar などが使うフック
export function useTree() {
    return useContext(TreeContext)
}
