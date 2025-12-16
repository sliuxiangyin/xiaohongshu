import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function LifecyclePage() {
  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div className="grid auto-rows-min gap-4">
        <Card>
          <CardHeader>
            <CardTitle>生命周期管理</CardTitle>
            <CardDescription>
              管理项目的生命周期阶段和状态
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p>这是生命周期管理页面的内容。</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}