import { AppSidebar } from "@/components/app-sidebar"
import { SiteHeader } from "@/components/site-header"
import {
  SidebarInset,
  SidebarProvider,
} from "@/components/ui/sidebar"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import HomePage from "./app/home/page"
import DashboardPage from "./app/dashboard/page"
import LifecyclePage from "./app/lifecycle/page"
import NotFoundPage from "./app/not-found/page"
import StartPage from "./app/start/page"

export default function App() {
  return (
    <Router>
      <Routes>
        {/* 独立的启动页面 */}
        <Route path="/" element={<StartPage />} />
        
        {/* 带侧边栏的主应用布局 */}
        <Route path="/*" element={
          <SidebarProvider
            style={
              {
                "--sidebar-width": "calc(var(--spacing) * 72)",
                "--header-height": "calc(var(--spacing) * 12)",
              } as React.CSSProperties
            }
          >
            <AppSidebar variant="inset" />
            <SidebarInset>
              <SiteHeader />
              <Routes>
                <Route path="/home" element={<HomePage />} />
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/lifecycle" element={<LifecyclePage />} />
              </Routes>
            </SidebarInset>
          </SidebarProvider>
        } />
        {/* 404 页面 */}
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Router>
  )
}