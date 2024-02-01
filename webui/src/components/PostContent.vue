<script>
    import { getImgUrl } from '../functions/getImgUrl';
    export default {
        props: ['post', 'containerClass', 'wrapperClass', 'likePhotoPath', 'username'],
        emits: ["change-class", "change-like", "delete-post"],
        methods: {
            getImgUrl,
        }
    }

</script>

<template>

    <div :class="containerClass">
        <div :class="wrapperClass">
            <img v-if="wrapperClass == 'post-popup'" @click="this.$emit('changeClass')" class="exit" src="@/assets/close.png" style="width: 30px; height: 30px;">
            <div class="post-header">
                <img class="author-img" :src="this.getImgUrl('profile_pics/' + post.Author + '.png')">
                <router-link :to="{ name: 'profile', params: {username: post.Author }}">
                    {{ post.Author }}
                </router-link>
                <span> {{ post.CreationDatetime }}</span>
                <img title="Delete post" @click="this.$emit('delete-post')" class="delete-post" v-if="post.Author == username && containerClass == 'post-container'" src="@/assets/x-red.png" >
            </div>
            <div class="post-body">
                <img class="post-image" @click="this.$emit('change-class')" :src="this.getImgUrl(post.Photo)">
                <div class="post-like">
                        Likes: {{ post.Likes ? post.Likes.length : 0 }}
                        <img @click="this.$emit('change-like')" :src="likePhotoPath">
                </div> 
                <div>
                    <span class="post-body-title"> {{ post.Author }}: </span>
                    <span class="post-body-description"> {{ post.Description }} </span> 
                </div>
            </div> 
        </div>
    </div>

</template>

<style scoped>

.post-container {
    padding: 1%;
    margin: 1%;
    border: 2px solid black;
    border-radius: 15px;
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
}

.post-header a {
    font-size: 18px;
    text-decoration: none;
    color: black
}

.post-header {
    position: relative
}

.post-header a:hover {
    font-weight: bold;
}

.post-header .author-img {
    margin-bottom: 0.5%;
    border-radius: 15px;
    width: 24px; 
    height: 24px;
}

.post-header .delete-post {
    border-radius: 15px;
    width: 24px; 
    height: 24px;
    position: absolute;
    top: 0;
    right: 0;
}

.post-header span {
    padding-left: 1%
}


.post-body-title {
    font-weight: bold;
    font-size: 15px;
}

.post-body-description{
    font-size: 15px;
}

.post-body .post-image {
    display: block; 
    max-width:100%; 
    max-height:200px; 
    width: auto; 
    height: auto;
    border-radius: 15px;
    box-shadow: rgba(0, 0, 0, 0.35) 0px 5px 15px;
}

.post-popup .post-image {
    display: block; 
    max-width: 100%; 
    max-height: 100%;
    width: auto; 
    height: auto;
    border-radius: 15px;
}

.post-like {
    padding-top: 1%;
    font-size: 18px;
}

.post-like img {
    width: 24px;
    height: 24px;
    margin-bottom: 4px;
}

.post-like img:hover {
    border: 1px solid black;
}

.post-overlay {

position: fixed;
top: 0;
bottom: 0;
left: 0;
right: 0;
background: black;
overflow:auto;
visibility: visible;
opacity: 1;
z-index: 1;
  
}

.post-popup {

margin: 70px auto;
padding: 20px;
background: #fff;
border-radius: 5px;
width: 40%;
position: relative;
background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);

}

</style>