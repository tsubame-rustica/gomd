import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { TreeProvider } from './TreeContext'
import Header from './Header'
import Home from './pages/Home'
import Post from './pages/Post'

function App() {
    return (
        <BrowserRouter>
            <TreeProvider>
                <div className="flex flex-col h-svh overflow-hidden">
                    <Header />
                    <div className="flex flex-row flex-1 overflow-hidden">
                        <Routes>
                            <Route path="/" element={<Home />} />
                            {/* URLPath が /git/default/default.md のような形式になるため /* で受ける */}
                            <Route path="/*" element={<Post />} />
                        </Routes>
                    </div>
                </div>
            </TreeProvider>
        </BrowserRouter>
    )
}

export default App