import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"

export default function NotFoundPage() {
  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div className="grid auto-rows-min gap-4">
        <Card>
          <CardHeader>
            <CardTitle>页面未找到</CardTitle>
            <CardDescription>
              抱歉，您访问的页面不存在
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="mb-4">您访问的页面可能已被移除或地址输入有误。</p>
            <Button asChild>
              <Link to="/">返回首页</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}