import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Header from './Header'
import Home from './pages/Home'
import Post from './pages/Post'


function App() {
    return (
        <BrowserRouter>
            <div className="flex flex-col h-svh overflow-hidden">
                <Header />
                <div className="flex flex-row flex-1 overflow-hidden">
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/:category/:title" element={<Post />} />
                    </Routes>
                </div>
            </div>
        </BrowserRouter>
    )
}

export default App