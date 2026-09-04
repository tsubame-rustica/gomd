import { Link } from 'react-router-dom'
import { useTree } from './TreeContext'
import type { DocumentNode } from './callApi'
import SearchBox from './SearchBox'

function Header() {
    return (
        <header className="shrink-0 w-full">
            <h1 className="text-3xl font-bold p-4">
                <Link to="/" className="text-neutral-800">
                    ドキュメント置き場
                </Link>
            </h1>
            <CategoryList />
        </header>
    )
}


function findFirstFile(node: DocumentNode): DocumentNode | undefined {
    if (node.isFile) return node
    for (const child of node.children) {
        const found = findFirstFile(child)
        if (found) return found
    }
    return undefined
}

function CategoryList() {
    const { tree, loading } = useTree()

    if (loading) return <nav className="px-6 py-2 bg-neutral-800" />

    // ツリーのルート直下のディレクトリ（カテゴリ）を表示
    const categories = tree?.children.filter(node => !node.isFile) ?? []

    return (
        <div className="flex flex-row items-center px-6 py-2 bg-neutral-800">
            <ul className="flex flex-row gap-4">
                <li key="home">
                    <Link to="/" className="text-white hover:underline">Home</Link>
                </li>
                {categories.map(category => (
                    <li key={category.urlPath}>
                        {/* order が最小のファイルにリンク */}
                        <Link
                            to={findFirstFile(category)?.urlPath ?? category.urlPath}
                            className="text-white hover:underline"
                        >
                            {category.displayName}
                        </Link>
                    </li>
                ))}
            </ul>
            <SearchBox />
        </div>
    )
}

export default Header