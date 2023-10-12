'use client'
import { useRef, useEffect } from "react"
import { useMap } from "../hooks/useMap"
import useSWR from "swr"
import { fetcher } from "../utils/http"
import { Route } from "../utils/model"
import { socket } from "../utils/socket-io"
import Grid2 from "@mui/material/Unstable_Grid2/Grid2"
import { Button, NativeSelect, Typography } from "@mui/material"

export function DriverPage() {
    const nestBaseUrl = process.env.NEXT_PUBLIC_NEST_BASE_URL
    const mapContainerRef = useRef<HTMLDivElement>(null)
    const map = useMap(mapContainerRef)
    const {data: routes, error, isLoading} = useSWR<Route[]>(`${nestBaseUrl}/routes`, fetcher, {fallbackData: []})
    
    useEffect(() => {
        socket.connect()
        return () => {
            socket.disconnect()
        }
    }, [])

    async function startRoute() {
        const routeId = (document.getElementById('routes') as HTMLSelectElement).value
        const response = await fetch(`${nestBaseUrl}/routes/${routeId}`)
        const route: Route = await response.json()
        map?.removeAllRoutes()
        await map?.addRouteWithIcons({
            routeId: routeId,
            startMarkerOptions: {
                position: route.directions.routes[0].legs[0].start_location,
            },
            endMarkerOptions: {
                position: route.directions.routes[0].legs[0].end_location,
            },
            carMarkerOptions: {
                position: route.directions.routes[0].legs[0].start_location,
            }
        })
        const {steps} = route.directions.routes[0].legs[0]
        for (const step of steps) {
            await sleep(1000)
            map?.moveCar(routeId, step.start_location)
            emitNewPoint(routeId, step)

            await sleep(1000)
            map?.moveCar(routeId, step.end_location)
            emitNewPoint(routeId, step)
        }
    }

    function emitNewPoint(routeId: string, step: any) {
        socket.emit('new-points', {
            route_id: routeId,
            lat: step.end_location.lat,
            lng: step.end_location.lng
        })
    }

    return (
        <Grid2 container sx={{ display: 'flex', flex: 1 }}>
            <Grid2 xs={4} px={2}>
                <Typography variant="h4">Minha viagem</Typography>
                <div style={{display:'flex', flexDirection: 'column'}}>
                    <NativeSelect id="routes">
                        {isLoading && <option>Carregando rotas...</option>}
                        {routes!.map((route) => (<option key={route.id} value={route.id}>{route.name}</option>))}
                    </NativeSelect>
                    <Button variant="contained" onClick={startRoute} fullWidth>Iniciar viagem</Button>
                </div>
            </Grid2>
            <Grid2 id="map" xs={8} ref={mapContainerRef}></Grid2>
        </Grid2>
    )
}

export default DriverPage

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))