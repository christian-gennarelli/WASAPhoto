<script>
    import Header from '../components/Header.vue'
    import Post from '../components/Post.vue'
    export default {
        data() {
            return{ 
                user: { // User info
                    token: '',
                    username: '',
                    name: '',
                    birthdate: '',
                    profilePic: '',
                },
                posts: [],
                postsCount: 0,
            }
        },
        methods: {
            getUserStream(){
                this.$axios.get(
                    '/users/' + this.user.username + '/stream',
                    { 
                        headers: {
                            'Authorization': this.user.token
                        }
                    }
                ).then((res) => {
                    this.posts = res.data
                    console.log(res.data)
                }).catch((e) => {
                    alert(e.data.ErrorCode + " " + e.data.Description)
                })
            },
        },
        components: {
            Header: Header,
            Post: Post
        },
        created() {

            // Username and token
            this.user.username = localStorage.getItem('username')
            this.user.token = localStorage.getItem('token')
            this.user.profilePic = localStorage.getItem('profilePic')
            this.user.birthdate = localStorage.getItem('birthdate')
            this.user.name = localStorage.getItem('name')

            console.log(this.user.token)

            // Stream
            this.getUserStream();
        }
    }
</script>

<template>
    <Header></Header>
    <p v-if="!posts"> No posts to see yet... start following!</p>
    <div>
        <Post 
            v-if="posts"
            v-for="post in posts" 
            :post="post" 
        ></Post>
    </div>
</template>