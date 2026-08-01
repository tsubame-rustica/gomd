import { useEffect, useState } from 'react'

interface Post {
    title: string;
    href: string;
    key: string;
}

interface Category {
    title: string;
    href: string;
    key: string;
}

function fetchPostsByCategory(category: string): { posts: Post[] } {
    const [posts, setPosts] = useState<Post[]>([]);

    useEffect(() => {
     fetch(`/api/getPosts/${category}`)
        .then(res => res.json())
        .then(data => {
            if (data.result) {
                setPosts(data.result);
            } else {
                console.error("No posts found for category:", category);
            }
        })
        .catch(err => {
            console.error("Fetch error:", err);
        });
    }, [posts]);

    return { posts };
}

function fetchCategories(): { categories: Category[] } {
    const [categories, setCategories] = useState<Category[]>([]);

    useEffect(() => {
     fetch(`/api/getCategories`)
        .then(res => res.json())
        .then(data => {
            if (data.result) {
                setCategories(data.result);
            } else {
                console.error("No categories found");
            }
        })
        .catch(err => {
            console.error("Fetch error:", err);
        });
    }, [categories]);

    return { categories };
}

export { fetchPostsByCategory, fetchCategories }