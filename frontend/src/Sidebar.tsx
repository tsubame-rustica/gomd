import { useRef, useState, useEffect, useCallback } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useTree } from './TreeContext'
import type { DocumentNode } from './callApi'

const MIN_WIDTH = 160
const MAX_WIDTH = 480
const DEFAULT_WIDTH = 256  // w-64 相当

function Sidebar() {
    const { tree } = useTree()
    const location = useLocation()
    const [width, setWidth] = useState(DEFAULT_WIDTH)
    const isDragging = useRef(false)
    const startX = useRef(0)
    const startWidth = useRef(0)

    const onMouseDown = useCallback((e: React.MouseEvent) => {
        isDragging.current = true
        startX.current = e.clientX
        startWidth.current = width
        // ドラッグ中にテキスト選択・カーソル変化を防ぐ
        document.body.style.userSelect = 'none'
        document.body.style.cursor = 'col-resize'
        e.preventDefault()
    }, [width])

    useEffect(() => {
        const onMouseMove = (e: MouseEvent) => {
            if (!isDragging.current) return
            const delta = e.clientX - startX.current
            const next = Math.min(MAX_WIDTH, Math.max(MIN_WIDTH, startWidth.current + delta))
            setWidth(next)
        }
        const onMouseUp = () => {
            if (!isDragging.current) return
            isDragging.current = false
            document.body.style.userSelect = ''
            document.body.style.cursor = ''
        }
        document.addEventListener('mousemove', onMouseMove)
        document.addEventListener('mouseup', onMouseUp)
        return () => {
            document.removeEventListener('mousemove', onMouseMove)
            document.removeEventListener('mouseup', onMouseUp)
        }
    }, [])

    // 現在のURLに一致するカテゴリ（ディレクトリ）を特定
    const currentCategory = tree?.children.find(
        node => !node.isFile && location.pathname.startsWith(node.urlPath)
    )
    const nodes = currentCategory?.children ?? []

    return (
        // relative で囲んでハンドルを絶対配置
        <div className="relative shrink-0 flex" style={{ width }}>
            {/* サイドバー本体 */}
            <div className="flex flex-col p-4 w-full bg-gray-100 overflow-y-auto">
                <ul className="flex flex-col">
                    {nodes.map(node => <NodeItem key={node.urlPath} node={node} />)}
                </ul>
            </div>

            {/* ドラッグハンドル（右端の細い帯） */}
            <div
                onMouseDown={onMouseDown}
                className="absolute right-0 top-0 h-full w-1 cursor-col-resize hover:bg-blue-400 transition-colors duration-150 z-10"
                title="ドラッグで幅を変更"
            />
        </div>
    )
}

function NodeItem({ node }: { node: DocumentNode }) {
    if (node.isFile) {
        return (
            <li>
                <Link to={node.urlPath} className="block text-neutral-800 p-2 hover:underline">
                    {node.displayName}
                </Link>
            </li>
        )
    }

    // サブディレクトリ: 折りたたみなしでラベルと子要素を表示
    return (
        <li>
            <ul>
                {node.children.map(child => <NodeItem key={child.urlPath} node={child} />)}
            </ul>
        </li>
    )
}

export default Sidebar