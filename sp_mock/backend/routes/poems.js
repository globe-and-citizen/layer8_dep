const express = require('express')
const cors = require('cors')
const Layer8 = require("../../../middleware/dist/loadWASM.js")
require('dotenv').config()

// INIT
var router = express.Router();

// MIDDLEWARE
router.use(cors())
router.use(express.json())
router.use(Layer8)

// Sample data (replace with actual data or database queries)
const poems = [
    {
        id: 1,
        title: 'The Raven',
        author: 'Edgar Allan Poe',
        content: `Once upon a midnight dreary, while I pondered, weak and weary,
      Over many a quaint and curious volume of forgotten lore—
      While I nodded, nearly napping, suddenly there came a tapping,
      As of some one gently rapping, rapping at my chamber door.
      "'Tis some visitor," I muttered, "tapping at my chamber door—
      Only this and nothing more."`
    },
    {
        id: 2,
        title: 'The Road Not Taken',
        author: 'Robert Frost',
        content: `Two roads diverged in a yellow wood,
      And sorry I could not travel both
      And be one traveler, long I stood
      And looked down one as far as I could
      To where it bent in the undergrowth;`
    },
    {
        id: 3,
        title: 'Ozymandias',
        author: 'Percy Bysshe Shelley',
        content: `I met a traveller from an antique land
      Who said: Two vast and trunkless legs of stone
      Stand in the desart. Near them, on the sand,
      Half sunk, a shattered visage lies, whose frown,`
    },
    {
        id: 4,
        title: 'The Love Song of J. Alfred Prufrock',
        author: 'T.S. Eliot',
        content: `Let us go then, you and I,
      When the evening is spread out against the sky
      Like a patient etherized upon a table;
      Let us go, through certain half-deserted streets,
      The muttering retreats`
    },
    {
        id: 5,
        title: 'Stopping by Woods on a Snowy Evening',
        author: 'Robert Frost',
        content: `Whose woods these are I think I know.
      His house is in the village, though;
      He will not see me stopping here
      To watch his woods fill up with snow.`
    }
];

router.get("/:id", (req, res) => {
    const poemId = parseInt(req.params.id);

    // Find the poem in the sample data (replace with your database query)
    const poem = poems.find((p) => p.id === poemId);

    if (!poem) {
        // Return a 404 Not Found response if the poem is not found
        return res.status(404).json({ error: 'Poem not found' });
    }

    // Return the found poem as JSON
    res.json(poem);
})

module.exports = router;
