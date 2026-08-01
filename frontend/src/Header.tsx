import { Link } from 'react-router-dom';

import { fetchCategories } from './callApi';


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

function CategoryList() {

    const { categories } = fetchCategories();

    const categoryList = categories.map((category) => (
        <li key={category.key}>
            <Link to={category.href} className="text-white hover:underline">
                {category.title}
            </Link>
        </li>
    ));

    return (
        <ul className="flex flex-row gap-4 px-6 py-2 bg-neutral-800">
            <li key="0">
                <Link to="/" className="text-white hover:underline">
                    Home
                </Link>
            </li>
            {categoryList}
        </ul>
    )
}


export default Header