twitter-thing back
---

## Setup
```
govendor fetch github.com/gin-gonic/gin@v1.2
govendor fetch github.com/gin-contrib/sessions
govendor fetch github.com/gin-contrib/cors
govendor fetch github.com/stretchr/testify/suite
npm i
```

## Develop
```
npm start
```
this is supposed to be keep building and restart local hosting whenever new build is available just like how node project does, using `nodemon`

but currently `nodemon` misteriously crashes when bulid fails, which is undesired behavior. Need to fix sometime...