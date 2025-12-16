import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ChartAreaInteractive } from "@/components/chart-area-interactive"

export default function DashboardPage() {
  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div className="grid auto-rows-min gap-4 md:grid-cols-4">
        <Card>
          <CardHeader>
            <CardTitle>用户统计</CardTitle>
            <CardDescription>总用户数</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold">1,234</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle>文章数量</CardTitle>
            <CardDescription>发布的文章</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold">567</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle>评论数</CardTitle>
            <CardDescription>用户评论</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold">8,901</p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader>
            <CardTitle>点赞数</CardTitle>
            <CardDescription>获得的点赞</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold">12,345</p>
          </CardContent>
        </Card>
      </div>
      
      <div className="grid auto-rows-min gap-4">
        <Card>
          <CardHeader>
            <CardTitle>数据趋势</CardTitle>
            <CardDescription>最近30天的数据变化</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartAreaInteractive />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}