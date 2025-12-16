import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Link, useNavigate } from "react-router-dom"
import { EventsEmit, EventsOff, EventsOn, LogPrint, } from "../../../wailsjs/runtime"
import { Spinner } from "@/components/ui/spinner"

export default function StartPage() {
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(true)
  const [errors, setErrors] = useState<string[]>([])
  useEffect(() => {
      LogPrint("发送开始准备信号 mounted")
      EventsEmit('startReady')
    // 监听初始化完成事件
    const cleanup = EventsOn("initialization-complete", (errorMessages: string[]) => {
      LogPrint("Received initialization-complete event")
      // 发送确认接收信号给后端
      EventsEmit('initialization-received')
      
      // 解析后端发送的错误信息
      if (errorMessages.length > 0) {
        // 如果有错误，设置错误状态但不跳转
        setErrors(errorMessages)
        setIsLoading(false)
        // 取消事件监听，避免重复处理
        EventsOff("initialization-complete")
      } else {
        // 如果没有错误，自动跳转到首页
        setIsLoading(false)
        // 取消事件监听，避免重复处理
        EventsOff("initialization-complete")
        navigate("/home")
      }
    })
    // 发送开始准备信号

    // 组件卸载时清理事件监听器
    return () => {
      cleanup()
    }
  }, [navigate])

  if (isLoading) {
    return (
      <div className="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
        <div className="flex flex-col items-center gap-4">
          <Spinner className="size-8 text-primary" />
          <p className="text-lg text-muted-foreground">加载中...</p>
        </div>
      </div>
    )
  }

  // 如果有错误，显示错误面板
  if (errors.length > 0) {
    return (
      <div className="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 p-4">
        <div className="flex flex-col items-center gap-6 max-w-2xl w-full">
          <Card className="w-full shadow-xl">
            <CardHeader className="text-center">
              <CardTitle className="text-3xl font-bold">初始化出现问题</CardTitle>
              <CardDescription className="text-lg">
                应用在初始化过程中遇到了一些问题
              </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-4">
              <div className="rounded-lg border border-red-200 bg-red-50 p-4 dark:border-red-800 dark:bg-red-950">
                <h3 className="font-medium text-red-800 dark:text-red-200">错误详情</h3>
                <p className="text-sm text-red-700 dark:text-red-300">
                  以下是初始化过程中遇到的错误：
                </p>
              </div>
              
              <div className="space-y-2 max-h-60 overflow-y-auto">
                {errors.map((error, index) => (
                  <div 
                    key={index} 
                    className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm dark:border-red-800 dark:bg-red-950"
                  >
                    <p className="text-red-800 dark:text-red-200 break-words">
                      {index + 1}. {error}
                    </p>
                  </div>
                ))}
              </div>
              
              <div className="flex flex-col sm:flex-row gap-3 pt-4">
                <Button 
                  onClick={() => window.location.reload()} 
                  className="w-full"
                >
                  重新尝试
                </Button>
                <Button 
                  variant="outline" 
                  asChild 
                  className="w-full"
                >
                  <Link to="/home">
                    忽略并继续
                  </Link>
                </Button>
              </div>
            </CardContent>
          </Card>
          
          <div className="text-center text-sm text-muted-foreground">
            <p>© 2023 小红书应用. 保留所有权利.</p>
          </div>
        </div>
      </div>
    )
  }

  // 正常情况下的欢迎页面
  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
      <div className="flex flex-col items-center gap-8 p-8 max-w-md w-full">
        <Card className="w-full shadow-xl">
          <CardHeader className="text-center">
            <CardTitle className="text-3xl font-bold">欢迎使用小红书应用</CardTitle>
            <CardDescription className="text-lg">
              您的社交媒体管理和分析平台
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col items-center gap-6">
            <p className="text-center text-muted-foreground">
              点击下方按钮开始探索应用的所有功能
            </p>
            <Button asChild size="lg" className="w-full">
              <Link to="/home">
                开始使用
              </Link>
            </Button>
          </CardContent>
        </Card>
        
        <div className="text-center text-sm text-muted-foreground">
          <p>© 2023 小红书应用. 保留所有权利.</p>
        </div>
      </div>
    </div>
  )
}