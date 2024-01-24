<script>
    import { getImgUrl } from '../functions/getImgUrl'
    export default {
        props: ['post'],
        methods: {
            getImgUrl
        }
    }
</script>

<template>
    <div>
        <div> 
            <img :src="this.getImgUrl('profile_pics/' + post.Author + '.png')" style="width: 16px; height: 16px"> 
            <router-link :to="{ name: 'profile', params: {username: post.Author }}">
                {{ post.Author }}
            </router-link>
        </div>
        <span>{{ post.CreationDatetime }}</span>
        <div>
            <img :src="this.getImgUrl(post.Photo)" style="width: 512px; height: 512px; border-radius: 20px; border: 2px solid black" >
        </div>
        <div> <span> Likes: {{ post.Likes }} </span> </div>
        <div> <span style="font-weight: bold;"> {{ post.Author }} </span> {{ post.Description }} </div>
    </div>
</template>
