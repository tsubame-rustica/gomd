import { Link } from 'react-router-dom'

import { fetchPostsByCategory } from './callApi'


function Sidebar({ category }: { category: string }) {
    return (
        <div className="flex flex-col p-4 w-64 shrink-0 bg-gray-100 overflow-y-auto">
            <ContentList category={category} />
        </div>
    )
}

function ContentList({ category }: { category: string }) {
    const { posts } = fetchPostsByCategory(category);

    const contentList = posts.map(post => (
        <li key={post.key}>
            <Link to={post.href} className="block text-neutral-800 p-2 hover:underline">
                {post.title}
            </Link>
        </li>
    ));

    return (
        <ul className="flex flex-col">
            {contentList}
        </ul>
    )
}

export default Sidebar