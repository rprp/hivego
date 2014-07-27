Spine = require('spineify')
Raphael = require('raphaelify')
Eve = require('eve')
$       = Spine.$

class Style

    @color=['#FF8C00', '#008000', '#2F4F4F', '#DA70D6', '#0000FF', '#8A2BE2', '#6495ED', '#B8860B', '#FF4500', '#AFEEEE', '#DB7093',
        '#CD853F', '#FFC0CB', '#B0E0E6', '#BC8F8F', '#4169E1', '#8B4513', '#00FFFF', '#00BFFF', '#008B8B',
        '#ADFF2F', '#4B0082', '#F0E68C', '#7CFC00', '#7FFF00', '#DEB887', '#98FB98', '#FFD700', '#5F9EA0', '#D2691E', '#A9A9A9',
        '#8B008B', '#556B2F', '#9932CC', '#8FBC8B', '#483D8B', '#00CED1', '#9400D3', '#FF69B4', '#228B22', '#1E90FF', '#FF00FF',
        '#FFB6C1', '#FFA07A', '#20B2AA', '#87CEFA', '#00FF00', '#B0C4DE', '#FF00FF', '#32CD32', '#0000CD', '#66CDAA', '#BA55D3',
        '#9370DB', '#3CB371', '#7B68EE', '#00FA9A', '#48D1CC', '#C71585', '#191970', '#000080', '#808000', '#6B8E23', '#FFA500',
        '#F4A460', '#2E8B57', '#A0522D', '#87CEEB', '#6A5ACD', '#708090', '#00FF7F', '#4682B4', '#D2B48C', '#008080', '#40E0D0',
         '#006400', '#BDB76B','#EE82EE', '#F5DEB3', '#FFFF00', '#9ACD32']

    [@sopt, @eopt] = [Raphael.animation({"fill-opacity": .2}, 1000), Raphael.animation({"fill-opacity": .5}, 1000)]

    [@csopt, @ceopt] = [Raphael.animation({"fill-opacity": .2}, 2000, -> @.animate(@ceopt)), Raphael.animation({"fill-opacity": .5}, 2000, -> @.animate(@csopt))]

    @slider = {fill: "#333", "fill-opacity": 0.3, "stroke-width": 2, "stroke-opacity": 0.1}
    
    @font = "Heiti, '黑体', 'Microsoft YaHei', '微软雅黑', SimSun, '宋体', '华文细黑', Helvetica, Tahoma, Arial, STXihei, sans-serif"
    @fontStyle = {fill: "#333", "font-family":@font, "text-anchor": "start", stroke: "none", "font-size": 18, "fill-opacity": 1, "stroke-width": 1}
    @jobFontStyle = {"font-family":@font , "font-size": 18, "stroke-opacity":1, "fill-opacity": 1, "stroke-width": 0}
    @jobcirStyle = {"fill-opacity": 0.2, "stroke-width": 1, cursor: "hand"}
    @jobrectStyle = {"fill-opacity": 0.1, "stroke-width": 0}
    @titlerectStyle = {fill: "#31708f", stroke: "#31708f", "fill-opacity": 0.05, "stroke-width": 0, cursor: "hand"}

module.exports = Style
