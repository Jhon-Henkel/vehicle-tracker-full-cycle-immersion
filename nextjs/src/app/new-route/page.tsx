'use client'
import type { FindPlaceFromTextResponseData } from "@googlemaps/google-maps-services-js"
import { FormEvent } from "react"

export function NewRoutePage() {    
    async function searchPlaces(event: FormEvent) {
        event.preventDefault()
        const source = (document.getElementById('source') as HTMLInputElement).value
        const destination = (document.getElementById('destination') as HTMLInputElement).value
        const [sourceResponse, destinationResponse] = await Promise.all([
            fetch(`http://localhost:3000/places?text=${source}`),
            fetch(`http://localhost:3000/places?text=${destination}`)
        ])
        const [sourcePlace, destinationPlace]: FindPlaceFromTextResponseData[] = await Promise.all([
            sourceResponse.json(),
            destinationResponse.json()
        ])
        if (sourcePlace.status !== 'OK') {
            alert('Origem não encontrada')
            return
        }
        if (destinationPlace.status !== 'OK') {
            alert('Destino não encontrada')
            return
        }
        const placeSourceId = sourcePlace.candidates[0].place_id
        const placeDestinationId = destinationPlace.candidates[0].place_id
        const directionsResponse = await fetch(
            `http://localhost:3000/directions?originId=${placeSourceId}&destinationId=${placeDestinationId}`
        )
        const directionsData = await directionsResponse.json()
    }
    return (
        <div>
            <h1>Nova Rota</h1>
            <form style={{display:'flex', flexDirection: 'column'}} onSubmit={searchPlaces}>
                <div>
                    <input id="source" type="text" placeholder="Origem" />
                </div>
                <div>
                    <input id="destination" type="text" placeholder="Destino" />
                </div>
                <button type="submit">Pesquisar</button>
            </form>
        </div>
    )
}

export default NewRoutePage