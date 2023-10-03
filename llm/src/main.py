from llm.Privatellm import PrivateLLM

if __name__ == "__main__":
    l = PrivateLLM()

    prompt = "A Nova Zelândia, uma fascinante nação insular no sudoeste do Oceano Pacífico, oferece uma combinação única de beleza natural deslumbrante e rica herança cultural. Suas paisagens dramáticas variam de montanhas cobertas de neve e florestas exuberantes a praias intocadas e fiordes. A cultura Maori, profundamente entrelaçada com a identidade da nação, adiciona um sabor distintivo com sua arte tradicional, dança e língua. A Nova Zelândia é famosa por suas aventuras ao ar livre, incluindo caminhadas, esqui e esportes aquáticos, tornando-se um paraíso para os entusiastas da natureza. O país também é celebrado por seu compromisso com a conservação ambiental, ostentando uma variedade de parques nacionais protegidos e santuários de vida selvagem. Além disso, seu povo caloroso e acolhedor, conhecido como Kiwis, contribui para a reputação do país de ser amigável e hospitaleiro. Em resumo, a Nova Zelândia é uma terra de paisagens deslumbrantes, diversidade cultural, emoções ao ar livre e uma forte dedicação à preservação de seus tesouros naturais."
    t="table"
    resp = l.PromptForBulletPoint(prompt, t)

    print(resp["history"])




